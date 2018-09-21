import { loadAll } from 'js-yaml';

export abstract class Kubeconfig {
  /**
   * @throws error if yaml is invalid
   */
  static parseYAML(yaml: string): any {
    if (yaml.length == 0) {
      return null;
    }

    let doc = loadAll(yaml);
    if (doc.length < 1) {
      throw new SyntaxError('Document does not look like a valid kubeconfig.');
    }

    let kubeconfig = doc[0];
    if (typeof kubeconfig.apiVersion !== 'string') {
      throw new SyntaxError('Document does not look like a valid kubeconfig.');
    }

    return kubeconfig;
  }
}
