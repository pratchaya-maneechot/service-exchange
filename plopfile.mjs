import { fileURLToPath } from "url";
import { dirname, join } from "path";

export default function (plop) {
  // --- Configuration ---
  const CONFIG = {
    apps: ["users", "tasks", "payment"], // Add your new apps here
    templatesPath: join(
      dirname(fileURLToPath(import.meta.url)),
      "tools/template/plop-go-templates"
    ),

    // Template file mappings
    templates: {
      id: join("go-ids", "id.hbs"),
      entity: join("go-domain", "entity.hbs"),
      repository_interface: join("go-domain", "repository_interface.hbs"),
      errors: join("go-domain", "errors.hbs"),
      event: join("go-domain", "event.hbs"),
      command: join("go-command", "command.hbs"),
      query: join("go-query", "query.hbs"),
      repository: join("go-repository", "repository.hbs"),
    },
  };

  // --- Reusable Prompts ---
  const prompts = {
    appName: {
      type: "list",
      name: "appName",
      message: "Which application (service) is this for?",
      choices: CONFIG.apps,
      default: CONFIG.apps[0],
    },

    moduleName: {
      type: "input",
      name: "moduleName",
      message:
        "What is the name of the module? (e.g., user, permission, product)",
      validate: (value) => (value.trim() ? true : "Module name is required"),
    },

    commandName: {
      type: "input",
      name: "commandName",
      message:
        "What is the name of the command? (e.g., RegisterUser, UpdateUserProfile)",
      validate: (value) => (value.trim() ? true : "Command name is required"),
    },

    queryName: {
      type: "input",
      name: "queryName",
      message:
        "What is the name of the query? (e.g., GetUserProfile, GetUserRoles)",
      validate: (value) => (value.trim() ? true : "Query name is required"),
    },
  };

  // --- Helpers ---
  plop.addHelper("snakeCase", function (text) {
    return text
      .replace(/([A-Z])/g, "_$1")
      .toLowerCase()
      .replace(/^_/, "");
  });

  // --- Action Builders ---
  const actionBuilders = {
    /**
     * Creates a file action with standardized paths and templates
     */
    createFile: (appName, relativePath, templateKey, data = {}) => ({
      type: "add",
      path: `apps/${appName}/${relativePath}`,
      templateFile: join(CONFIG.templatesPath, CONFIG.templates[templateKey]),
      data,
    }),

    /**
     * Creates domain layer files
     */
    createDomainFiles: (appName, moduleName) => [
      actionBuilders.createFile(
        appName,
        `internal/domain/shared/ids/{{snakeCase moduleName}}_id.go`,
        "id",
        { moduleName }
      ),
      actionBuilders.createFile(
        appName,
        `internal/domain/{{snakeCase moduleName}}/{{snakeCase moduleName}}.go`,
        "entity",
        { moduleName }
      ),
      actionBuilders.createFile(
        appName,
        `internal/domain/{{snakeCase moduleName}}/repository.go`,
        "repository_interface",
        { moduleName }
      ),
      actionBuilders.createFile(
        appName,
        `internal/domain/{{snakeCase moduleName}}/errors.go`,
        "errors",
        { moduleName }
      ),
      actionBuilders.createFile(
        appName,
        `internal/domain/{{snakeCase moduleName}}/event.go`,
        "event",
        { moduleName }
      ),
    ],

    /**
     * Creates CRUD command files
     */
    createCrudCommands: (appName, moduleName) => [
      actionBuilders.createFile(
        appName,
        `internal/app/command/create_{{snakeCase moduleName}}.go`,
        "command",
        {
          moduleName,
          commandName: `Create${plop.getHelper("pascalCase")(moduleName)}`,
        }
      ),
      actionBuilders.createFile(
        appName,
        `internal/app/command/update_{{snakeCase moduleName}}.go`,
        "command",
        {
          moduleName,
          commandName: `Update${plop.getHelper("pascalCase")(moduleName)}`,
        }
      ),
      actionBuilders.createFile(
        appName,
        `internal/app/command/delete_{{snakeCase moduleName}}.go`,
        "command",
        {
          moduleName,
          commandName: `Delete${plop.getHelper("pascalCase")(moduleName)}`,
        }
      ),
    ],

    /**
     * Creates standard query files
     */
    createStandardQueries: (appName, moduleName) => [
      actionBuilders.createFile(
        appName,
        `internal/app/query/get_{{snakeCase moduleName}}s.go`,
        "query",
        {
          moduleName,
          queryName: `Get${plop.getHelper("pascalCase")(moduleName)}s`,
        }
      ),
      actionBuilders.createFile(
        appName,
        `internal/app/query/get_{{snakeCase moduleName}}.go`,
        "query",
        {
          moduleName,
          queryName: `Get${plop.getHelper("pascalCase")(moduleName)}`,
        }
      ),
    ],

    /**
     * Creates repository implementation
     */
    createRepository: (appName, moduleName) => [
      actionBuilders.createFile(
        appName,
        `internal/infra/persistence/repositories/{{snakeCase moduleName}}_repository.go`,
        "repository",
        { moduleName }
      ),
    ],
  };

  // --- Generators ---

  /**
   * Complete module generator - creates all files for a new module
   */
  plop.setGenerator("module", {
    description:
      "Generates a complete internal module (Domain, Commands, Queries, Repository)",
    prompts: [prompts.appName, prompts.moduleName],
    actions: (data) => {
      const { appName, moduleName } = data;

      return [
        ...actionBuilders.createDomainFiles(appName, moduleName),
        ...actionBuilders.createCrudCommands(appName, moduleName),
        ...actionBuilders.createStandardQueries(appName, moduleName),
        ...actionBuilders.createRepository(appName, moduleName),
      ];
    },
  });

  /**
   * Single command generator
   */
  plop.setGenerator("command", {
    description: "Generates a new Go command handler",
    prompts: [prompts.appName, prompts.moduleName, prompts.commandName],
    actions: (data) => {
      const { appName, moduleName, commandName } = data;
      const snakeCommandName = plop.getHelper("snakeCase")(commandName);

      return [
        actionBuilders.createFile(
          appName,
          `internal/app/command/${snakeCommandName}.go`,
          "command",
          { moduleName, commandName }
        ),
      ];
    },
  });

  /**
   * Single query generator
   */
  plop.setGenerator("query", {
    description: "Generates a new Go query handler",
    prompts: [prompts.appName, prompts.moduleName, prompts.queryName],
    actions: (data) => {
      const { appName, moduleName, queryName } = data;
      const snakeQueryName = plop.getHelper("snakeCase")(queryName);

      return [
        actionBuilders.createFile(
          appName,
          `internal/app/query/get_${snakeQueryName}.go`, // Adjusted to include 'get_' prefix
          "query",
          { moduleName, queryName }
        ),
      ];
    },
  });

  /**
   * Repository implementation generator
   */
  plop.setGenerator("repository", {
    description: "Generates a new Go repository implementation",
    prompts: [prompts.appName, prompts.moduleName],
    actions: (data) =>
      actionBuilders.createRepository(data.appName, data.moduleName),
  });

  /**
   * Domain files generator
   */
  plop.setGenerator("domain", {
    description:
      "Generates core domain files (Entity, Errors, Event, Repository Interface, ID)",
    prompts: [prompts.appName, prompts.moduleName],
    actions: (data) =>
      actionBuilders.createDomainFiles(data.appName, data.moduleName),
  });

  /**
   * CRUD commands generator
   */
  plop.setGenerator("crud-commands", {
    description:
      "Generates standard CRUD command handlers (Create, Update, Delete)",
    prompts: [prompts.appName, prompts.moduleName],
    actions: (data) =>
      actionBuilders.createCrudCommands(data.appName, data.moduleName),
  });

  /**
   * Standard queries generator
   */
  plop.setGenerator("standard-queries", {
    description: "Generates standard query handlers (GetAll, GetById)",
    prompts: [prompts.appName, prompts.moduleName],
    actions: (data) =>
      actionBuilders.createStandardQueries(data.appName, data.moduleName),
  });
}