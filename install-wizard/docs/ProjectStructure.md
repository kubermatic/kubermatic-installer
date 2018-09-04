# Project Structure

## Folder Structure

| Location                    | Description                                                               |
|:----------------------------|:--------------------------------------------------------------------------|
| /                           | Contains the whole wizard ui project                                      |
| /docs                       | Contains docs for the ui project                                          |
| /src/styles.scss            | The global stylesheets, no component related things should be placed here |
| /src/assets/                | Contains binary assets as in: images, fonts, ...                          |
| /src/app/                   | Contains the root app component                                           |
| /src/app/services/          | Should contain all services                                               |
| /src/app/models/            | Should contain the shared models                                          |
| /src/app/interfaces/        | Should contain shared interfaces                                          |
| /src/app/components/wizard/ | Contains the main wizard component                                        |
| /src/app/components/steps/  | Should contain a component for each wizard step we have                   |
| /e2e/src/                   | Contains e2e tests for the wizard                                         |

Convention should be, that each component(html + scss + spec + ts) should be placed into an own directory.

## Architecture

The manifest should be part of the wizard component, the individual steps of it
should be bound the the corresponding step components.