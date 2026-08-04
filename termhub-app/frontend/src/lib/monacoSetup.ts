// Bundles Monaco locally instead of the @monaco-editor/react default of
// fetching it from a CDN at runtime — this app can be deployed on a private
// VPN/intranet with no outbound internet access, so the editor must work
// fully offline. Side-effect module: import it once before any <Editor/>
// mounts (done in EditorView.tsx, which is lazy-loaded so this cost is only
// paid once the user actually opens the editor tab).
import * as monaco from 'monaco-editor';
// monaco-editor's package "exports" map already prefixes esm/vs/, so
// subpaths here must NOT repeat it (e.g. 'monaco-editor/editor/editor.worker',
// not 'monaco-editor/esm/vs/editor/editor.worker').
import EditorWorker from 'monaco-editor/editor/editor.worker?worker';
import JsonWorker from 'monaco-editor/language/json/json.worker?worker';
import CssWorker from 'monaco-editor/language/css/css.worker?worker';
import HtmlWorker from 'monaco-editor/language/html/html.worker?worker';
import TsWorker from 'monaco-editor/language/typescript/ts.worker?worker';
import { loader } from '@monaco-editor/react';

self.MonacoEnvironment = {
  getWorker(_workerId: string, label: string) {
    switch (label) {
      case 'json':
        return new JsonWorker();
      case 'css':
      case 'scss':
      case 'less':
        return new CssWorker();
      case 'html':
      case 'handlebars':
      case 'razor':
        return new HtmlWorker();
      case 'typescript':
      case 'javascript':
        return new TsWorker();
      default:
        return new EditorWorker();
    }
  },
};

loader.config({ monaco });
