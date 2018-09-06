# Wizard UI

Use Docker to not inject your system with Node.js and all its madness.
The Makefile defines a couple of helpful tasks:

* `make dockerbuild` builds the Docker image and needs to be run first.
* `make npminstall` installs the npm dependencies.
* `make npmstart` starts the app; it uses the host network stack, so
  you should be able to access it on `http://127.0.0.1:4200` afterwards.
* `make shell` hands you a shell to work in.

The wizard is based on [angular-archwizard](https://github.com/madoar/angular-archwizard)
and works like this:

* We have a `WizardComponent` that represents the wizard itself.
* Each wizard step is its own component (inheriting from `WizardStep`
  to make things easier), named after their task (so they are not
  enumerated, but given proper names, in case we want to shuffle them
  around at some point).
* A single, simple `Manifest` class represents the data structure that
  all wizard steps work on. They are all given access to it by input
  injection. Right now the manifest only contains a proof-of-concept
  field for the chosen cloud provider.
* Cloud providers are hardcoded in the `config.ts`.
* Each wizard step defines an Angular Reactive Form in order to be able
  to have custom validation logic (outside of the HTML-provided attributes
  like `required`, `minlength` etc.).
* When entering stuff into the form, at first all the validators for the
  field will be called, then the state will be synced to the underlying
  model (using your syncManifestFromForm callback) and then the form-wide
  validation is called.

To add a new wizard step:

* `ng create component wizard-step-something-something-dark-side`.
* Add the `<app-wizard-step-something-something-dark-side>` to
  `app/wizard/wizard.component.html`.
* Hack your wizard step page, i.e. add a form, add cat gifs, whatever.
  It's probably a good idea to copy stuff from existing steps.
