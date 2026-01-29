"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const vscode = __importStar(require("vscode"));
const erstClient_1 = require("./erstClient");
const traceTreeView_1 = require("./traceTreeView");
function activate(context) {
    const client = new erstClient_1.ERSTClient('127.0.0.1', 8080);
    const traceDataProvider = new traceTreeView_1.TraceTreeDataProvider();
    // Register TreeView
    vscode.window.registerTreeDataProvider('erst-traces', traceDataProvider);
    // Register command: erst.triggerDebug
    let triggerDebugDisposable = vscode.commands.registerCommand('erst.triggerDebug', async () => {
        const hash = await vscode.window.showInputBox({
            prompt: 'Enter Transaction Hash to Debug',
            placeHolder: 'e.g., sample-tx-hash-1234'
        });
        if (hash) {
            try {
                await vscode.window.withProgress({
                    location: vscode.ProgressLocation.Notification,
                    title: "ERST: Debugging Transaction...",
                    cancellable: false
                }, async (progress) => {
                    await client.connect();
                    await client.debugTransaction(hash);
                    const trace = await client.getTrace(hash);
                    traceDataProvider.refresh(trace);
                });
                vscode.window.showInformationMessage(`Trace loaded for ${hash}`);
            }
            catch (err) {
                vscode.window.showErrorMessage(`ERST Error: ${err.message}`);
            }
        }
    });
    // Handle selecting a trace item
    let selectTraceStepDisposable = vscode.commands.registerCommand('erst.selectTraceStep', (item) => {
        const stepJson = JSON.stringify(item.step, null, 2);
        // Show in a virtual document or just a message for PoC
        vscode.workspace.openTextDocument({
            content: stepJson,
            language: 'json'
        }).then(doc => {
            vscode.window.showTextDocument(doc, vscode.ViewColumn.Beside);
        });
    });
    context.subscriptions.push(triggerDebugDisposable, selectTraceStepDisposable, client);
}
function deactivate() { }
//# sourceMappingURL=extension.js.map