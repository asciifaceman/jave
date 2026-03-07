const vscode = require("vscode");
const cp = require("child_process");
const fs = require("fs");
const path = require("path");

class LspClient {
  constructor(command, args, cwd) {
    this.command = command;
    this.args = args;
    this.cwd = cwd;
    this.proc = null;
    this.seq = 1;
    this.pending = new Map();
    this.buffer = Buffer.alloc(0);
    this.initialized = false;
  }

  start() {
    if (!this.command) {
      return Promise.reject(new Error("No LSP command configured"));
    }

    this.proc = cp.spawn(this.command, this.args, {
      cwd: this.cwd || undefined,
      stdio: ["pipe", "pipe", "pipe"],
    });

    this.proc.on("error", (err) => {
      for (const [, reject] of this.pending.values()) {
        reject(err);
      }
      this.pending.clear();
      this.initialized = false;
    });

    this.proc.stderr.on("data", () => {
      // Keep stderr drained so child process cannot block on full pipe.
    });

    this.proc.stdout.on("data", (chunk) => this.onData(chunk));
    this.proc.on("exit", () => {
      for (const [, reject] of this.pending.values()) {
        reject(new Error("javels exited"));
      }
      this.pending.clear();
      this.initialized = false;
    });

    const rootUri = vscode.workspace.workspaceFolders?.[0]?.uri?.toString() || null;
    return this.request("initialize", {
      processId: process.pid,
      rootUri,
      capabilities: {},
      clientInfo: { name: "vscode-jave", version: "0.1.0" },
    }).then(() => {
      this.notify("initialized", {});
      this.initialized = true;
    });
  }

  stop() {
    if (!this.proc) {
      return Promise.resolve();
    }
    return this.request("shutdown", {})
      .catch(() => undefined)
      .then(() => {
        this.notify("exit", {});
        this.proc.kill();
      });
  }

  request(method, params) {
    if (!this.proc || !this.proc.stdin.writable) {
      return Promise.reject(new Error("javels not running"));
    }
    const id = this.seq++;
    const payload = JSON.stringify({ jsonrpc: "2.0", id, method, params });
    this.write(payload);
    return new Promise((resolve, reject) => {
      this.pending.set(id, [resolve, reject]);
      setTimeout(() => {
        if (this.pending.has(id)) {
          this.pending.delete(id);
          reject(new Error(`javels timeout for ${method}`));
        }
      }, 5000);
    });
  }

  notify(method, params) {
    if (!this.proc || !this.proc.stdin.writable) {
      return;
    }
    const payload = JSON.stringify({ jsonrpc: "2.0", method, params });
    this.write(payload);
  }

  write(payload) {
    const header = `Content-Length: ${Buffer.byteLength(payload, "utf8")}\r\n\r\n`;
    this.proc.stdin.write(header + payload);
  }

  onData(chunk) {
    this.buffer = Buffer.concat([this.buffer, chunk]);
    while (true) {
      const marker = this.buffer.indexOf("\r\n\r\n");
      if (marker < 0) {
        return;
      }
      const header = this.buffer.slice(0, marker).toString("utf8");
      const m = /Content-Length:\s*(\d+)/i.exec(header);
      if (!m) {
        this.buffer = Buffer.alloc(0);
        return;
      }
      const len = parseInt(m[1], 10);
      const total = marker + 4 + len;
      if (this.buffer.length < total) {
        return;
      }
      const body = this.buffer.slice(marker + 4, total).toString("utf8");
      this.buffer = this.buffer.slice(total);
      this.handleMessage(body);
    }
  }

  handleMessage(body) {
    let msg;
    try {
      msg = JSON.parse(body);
    } catch {
      return;
    }
    if (msg.id === undefined || msg.id === null) {
      return;
    }
    const pending = this.pending.get(msg.id);
    if (!pending) {
      return;
    }
    this.pending.delete(msg.id);
    const [resolve, reject] = pending;
    if (msg.error) {
      reject(new Error(msg.error.message || "LSP error"));
      return;
    }
    resolve(msg.result);
  }
}

function workspaceCwd(config) {
  const explicit = config.get("languageServer.cwd");
  if (explicit && explicit.trim() !== "") {
    return explicit;
  }
  return vscode.workspace.workspaceFolders?.[0]?.uri?.fsPath;
}

function bundledServerPath(extensionPath) {
  const platform = process.platform;
  const arch = process.arch;
  let osName;
  if (platform === "win32") {
    osName = "windows";
  } else if (platform === "darwin") {
    osName = "darwin";
  } else if (platform === "linux") {
    osName = "linux";
  } else {
    return null;
  }

  let archName;
  if (arch === "x64") {
    archName = "amd64";
  } else if (arch === "arm64") {
    archName = "arm64";
  } else {
    return null;
  }

  const exe = osName === "windows" ? ".exe" : "";
  const candidate = path.join(extensionPath, "bin", `javels-${osName}-${archName}${exe}`);
  if (fs.existsSync(candidate)) {
    return candidate;
  }
  return null;
}

function languageServerLaunch(context, config) {
  const useBundled = config.get("languageServer.useBundled") !== false;
  if (useBundled) {
    const bundled = bundledServerPath(context.extensionPath);
    if (bundled) {
      return {
        command: bundled,
        args: [],
        cwd: workspaceCwd(config),
        mode: "bundled",
      };
    }
  }

  const command = config.get("languageServer.command") || "javels";
  const args = config.get("languageServer.args") || [];
  return {
    command,
    args,
    cwd: workspaceCwd(config),
    mode: "configured",
  };
}

function devGoRunLaunch(config) {
  const cwd = workspaceCwd(config);
  if (!cwd) {
    return null;
  }
  const candidate = path.join(cwd, "cmd", "javels", "main.go");
  if (!fs.existsSync(candidate)) {
    return null;
  }
  return {
    command: "go",
    args: ["run", "./cmd/javels"],
    cwd,
    mode: "dev-go-run",
  };
}

async function startClientWithFallback(context, cfg) {
  const primary = languageServerLaunch(context, cfg);
  const launches = [primary];

  if (primary.mode !== "configured") {
    launches.push({
      command: cfg.get("languageServer.command") || "javels",
      args: cfg.get("languageServer.args") || [],
      cwd: workspaceCwd(cfg),
      mode: "configured",
    });
  }

  const dev = devGoRunLaunch(cfg);
  if (dev) {
    launches.push(dev);
  }

  let lastErr = null;
  for (const launch of launches) {
    const client = new LspClient(launch.command, launch.args, launch.cwd);
    try {
      await client.start();
      return { client, launch };
    } catch (err) {
      lastErr = err;
    }
  }

  throw lastErr || new Error("Jave LSP launch failed");
}

function activate(context) {
  const cfg = vscode.workspace.getConfiguration("jave");
  let client = null;

  startClientWithFallback(context, cfg)
    .then((result) => {
      client = result.client;
      if (result.launch.mode !== "bundled") {
        vscode.window.setStatusBarMessage(
          `Jave LSP running via ${result.launch.mode} launcher`,
          4000
        );
      }
    })
    .catch((err) => {
      vscode.window.showWarningMessage(`Jave LSP failed to start: ${err.message}`);
    });

  const syncOpen = (doc) => {
    if (doc.languageId !== "jave") {
      return;
    }
    if (!client) {
      return;
    }
    client.notify("textDocument/didOpen", {
      textDocument: {
        uri: doc.uri.toString(),
        text: doc.getText(),
      },
    });
  };

  const syncChange = (event) => {
    if (event.document.languageId !== "jave") {
      return;
    }
    if (!client) {
      return;
    }
    client.notify("textDocument/didChange", {
      textDocument: {
        uri: event.document.uri.toString(),
      },
      contentChanges: [{ text: event.document.getText() }],
    });
  };

  for (const doc of vscode.workspace.textDocuments) {
    syncOpen(doc);
  }

  context.subscriptions.push(vscode.workspace.onDidOpenTextDocument(syncOpen));
  context.subscriptions.push(vscode.workspace.onDidChangeTextDocument(syncChange));

  const hoverProvider = vscode.languages.registerHoverProvider("jave", {
    async provideHover(document, position) {
      try {
        if (!client) {
          return null;
        }
        const result = await client.request("textDocument/hover", {
          textDocument: { uri: document.uri.toString() },
          position: { line: position.line, character: position.character },
        });
        const value = result?.contents?.value;
        if (!value) {
          return null;
        }
        return new vscode.Hover(new vscode.MarkdownString(value));
      } catch {
        return null;
      }
    },
  });

  const signatureProvider = vscode.languages.registerSignatureHelpProvider(
    "jave",
    {
      async provideSignatureHelp(document, position) {
        try {
          if (!client) {
            return null;
          }
          const result = await client.request("textDocument/signatureHelp", {
            textDocument: { uri: document.uri.toString() },
            position: { line: position.line, character: position.character },
          });
          if (!result || !Array.isArray(result.signatures) || result.signatures.length === 0) {
            return null;
          }
          const help = new vscode.SignatureHelp();
          help.activeSignature = result.activeSignature || 0;
          help.activeParameter = result.activeParameter || 0;
          help.signatures = result.signatures.map((sig) => {
            const s = new vscode.SignatureInformation(sig.label, sig.documentation?.value || "");
            s.parameters = (sig.parameters || []).map((p) => new vscode.ParameterInformation(p.label, p.documentation?.value || ""));
            return s;
          });
          return help;
        } catch {
          return null;
        }
      },
    },
    "<",
    "(",
    ","
  );

  context.subscriptions.push(hoverProvider, signatureProvider);
  context.subscriptions.push({
    dispose: () => {
      if (client) {
        return client.stop();
      }
      return Promise.resolve();
    },
  });
}

function deactivate() {
  // VS Code calls disposables in activate cleanup.
}

module.exports = {
  activate,
  deactivate,
};
