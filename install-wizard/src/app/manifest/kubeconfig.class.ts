import { loadAll } from 'js-yaml';

export abstract class Kubeconfig {
  /**
   * @throws error if yaml is invalid
   */
  static parseYAML(yaml: string): any {
    if (yaml.length === 0) {
      return null;
    }

    const doc = loadAll(yaml);
    if (doc.length < 1) {
      throw new SyntaxError('could not parse YAML');
    }

    const kubeconfig = doc[0];
    if (typeof kubeconfig.apiVersion !== 'string') {
      throw new SyntaxError('no apiVersion defined');
    }

    return kubeconfig;
  }

  static getContexts(kubeconfig: any): string[] {
    if (typeof kubeconfig.contexts === 'undefined' || !Array.isArray(kubeconfig.contexts)) {
      throw new SyntaxError('no contexts array defined');
    }

    const contexts = [];

    kubeconfig.contexts.forEach(context => {
      if (!contexts.includes(context.name)) {
        contexts.push(context.name);
      }
    });

    return contexts.sort();
  }

  /**
   * Checks if a context name is valid, i.e. can be used inside DNS
   * names to construct the full domain for a seed cluster.
   *
   * @param name string
   */
  static isValidContextName(name: string): boolean {
    if (name.length === 0 || name.length > 63) {
      return false;
    }

    if (!/^[a-z0-9-]+$/.test(name)) {
      return false;
    }

    if (name.startsWith('-') || name.endsWith('-')) {
      return false;
    }

    return true;
  }
}
