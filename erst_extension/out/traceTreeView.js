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
exports.TraceItem = exports.TraceTreeDataProvider = void 0;
const vscode = __importStar(require("vscode"));
class TraceTreeDataProvider {
    _onDidChangeTreeData = new vscode.EventEmitter();
    onDidChangeTreeData = this._onDidChangeTreeData.event;
    currentTrace;
    constructor() { }
    refresh(trace) {
        this.currentTrace = trace;
        this._onDidChangeTreeData.fire();
    }
    getTreeItem(element) {
        return element;
    }
    getChildren(element) {
        if (!this.currentTrace) {
            return Promise.resolve([]);
        }
        if (element) {
            // Further details if expanded, but for now we just show steps
            return Promise.resolve([]);
        }
        else {
            return Promise.resolve(this.currentTrace.states.map(step => new TraceItem(step)));
        }
    }
}
exports.TraceTreeDataProvider = TraceTreeDataProvider;
class TraceItem extends vscode.TreeItem {
    step;
    constructor(step) {
        super(`${step.step}: ${step.operation}${step.function ? ` (${step.function})` : ''}`, vscode.TreeItemCollapsibleState.None);
        this.step = step;
        this.tooltip = `${this.label}`;
        this.description = step.error ? `Error: ${step.error}` : '';
        this.contextValue = 'traceStep';
        if (step.error) {
            this.iconPath = new vscode.ThemeIcon('error', new vscode.ThemeColor('errorForeground'));
        }
        else {
            this.iconPath = new vscode.ThemeIcon('pass', new vscode.ThemeColor('debugIcon.startForeground'));
        }
    }
}
exports.TraceItem = TraceItem;
//# sourceMappingURL=traceTreeView.js.map