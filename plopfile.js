// plopfile.js
module.exports = function (plop) {
  // --- Helpers ---
  // Plop has default helpers for pascalCase, kebabCase, camelCase
  // We'll add snakeCase if it's not already there
  plop.addHelper('snakeCase', function (text) {
    return text.replace(/([A-Z])/g, '_$1').toLowerCase();
  });

  // --- Common Prompts ---
  const moduleNamePrompt = {
    type: 'input',
    name: 'moduleName',
    message: 'What is the name of the module? (e.g., user, permission, product)',
    validate: function (value) {
      if ((/.+/).test(value)) { return true; }
      return 'module name is required';
    }
  };

  const commandNamePrompt = {
    type: 'input',
    name: 'commandName',
    message: 'What is the name of the command? (e.g., RegisterUser, UpdateUserProfile)',
    validate: function (value) {
      if ((/.+/).test(value)) { return true; }
      return 'command name is required';
    }
  };

  const queryNamePrompt = {
    type: 'input',
    name: 'queryName',
    message: 'What is the name of the query? (e.g., GetUserProfile, GetUserRoles)',
    validate: function (value) {
      if ((/.+/).test(value)) { return true; }
      return 'query name is required';
    }
  };

  // --- Generators ---

  // --- Generator for a full Internal Module ---
  plop.setGenerator('internal-module', {
    description: 'Generates a new complete internal module (Domain, Commands, Queries, Repository)',
    prompts: [moduleNamePrompt],
    actions: function(data) {
      const actions = [];
      const moduleName = data.moduleName; // e.g., 'Permission'

      // 1. Domain Layer
      actions.push({
        type: 'add',
        path: 'apps/users/internal/domain/shared/ids/{{kebabCase moduleName}}_id.go',
        templateFile: 'tools/template/plop-go-templates/go-ids/id.hbs',
        data: { moduleName: moduleName }
      });
      actions.push({
        type: 'add',
        path: 'apps/users/internal/domain/{{kebabCase moduleName}}/{{kebabCase moduleName}}.go',
        templateFile: 'tools/template/plop-go-templates/go-domain/entity.hbs',
        data: { moduleName: moduleName }
      });
      actions.push({
        type: 'add',
        path: 'apps/users/internal/domain/{{kebabCase moduleName}}/repository.go',
        templateFile: 'tools/template/plop-go-templates/go-domain/repository_interface.hbs',
        data: { moduleName: moduleName }
      });
      actions.push({
        type: 'add',
        path: 'apps/users/internal/domain/{{kebabCase moduleName}}/errors.go',
        templateFile: 'tools/template/plop-go-templates/go-domain/errors.hbs',
        data: { moduleName: moduleName }
      });
      actions.push({
        type: 'add',
        path: 'apps/users/internal/domain/{{kebabCase moduleName}}/event.go',
        templateFile: 'tools/template/plop-go-templates/go-domain/event.hbs',
        data: { moduleName: moduleName }
      });

      // 2. Application Layer (Commands)
      actions.push({
        type: 'add',
        path: 'apps/users/internal/app/command/create_{{kebabCase moduleName}}.go',
        templateFile: 'tools/template/plop-go-templates/go-command/command.hbs',
        data: { moduleName: moduleName, commandName: `Create${moduleName}` }
      });
      actions.push({
        type: 'add',
        path: 'apps/users/internal/app/command/update_{{kebabCase moduleName}}.go',
        templateFile: 'tools/template/plop-go-templates/go-command/command.hbs',
        data: { moduleName: moduleName, commandName: `Update${moduleName}` }
      });
      actions.push({
        type: 'add',
        path: 'apps/users/internal/app/command/delete_{{kebabCase moduleName}}.go',
        templateFile: 'tools/template/plop-go-templates/go-command/command.hbs',
        data: { moduleName: moduleName, commandName: `Delete${moduleName}` }
      });

      // 3. Application Layer (Queries)
      actions.push({
        type: 'add',
        path: 'apps/users/internal/app/query/get_{{kebabCase moduleName}}s.go',
        templateFile: 'tools/template/plop-go-templates/go-query/query.hbs',
        data: { moduleName: moduleName, queryName: `Get${moduleName}s` }
      });
      actions.push({
        type: 'add',
        path: 'apps/users/internal/app/query/get_{{kebabCase moduleName}}.go',
        templateFile: 'tools/template/plop-go-templates/go-query/query.hbs',
        data: { moduleName: moduleName, queryName: `Get${moduleName}` }
      });

      // 4. Infrastructure Layer (Repository)
      actions.push({
        type: 'add',
        path: 'apps/users/internal/infra/persistence/repositories/{{kebabCase moduleName}}_repository.go',
        templateFile: 'tools/template/plop-go-templates/go-repository/repository.hbs',
        data: { moduleName: moduleName }
      });

      return actions;
    }
  });

  // --- Generator for a single Internal Command ---
  plop.setGenerator('internal-command', {
    description: 'Generates a new Go command handler',
    prompts: [
      moduleNamePrompt,
      commandNamePrompt
    ],
    actions: function(data) {
      // Create a new data object for the template, applying helpers for naming
      const templateData = {
        moduleName: data.moduleName,
        commandName: data.commandName,
        // Ensure kebabCase and pascalCase are applied to the commandName for the path and template
        kebabCommandName: plop.getHelper('kebabCase')(data.commandName),
        pascalCommandName: plop.getHelper('pascalCase')(data.commandName)
      };

      return [
        {
          type: 'add',
          path: `apps/users/internal/app/command/${templateData.kebabCommandName}.go`, // Use pre-transformed name
          templateFile: 'tools/template/plop-go-templates/go-command/command.hbs',
          data: templateData // Pass the pre-processed data to the template
        }
      ];
    }
  });

  // --- Generator for a single Internal Query ---
  plop.setGenerator('internal-query', {
    description: 'Generates a new Go query handler',
    prompts: [
      moduleNamePrompt,
      queryNamePrompt
    ],
    actions: function(data) {
      const templateData = {
        moduleName: data.moduleName,
        queryName: data.queryName,
        kebabQueryName: plop.getHelper('kebabCase')(data.queryName),
        pascalQueryName: plop.getHelper('pascalCase')(data.queryName)
      };

      return [
        {
          type: 'add',
          path: `apps/users/internal/app/query/get_${templateData.kebabQueryName}.go`, // Use pre-transformed name
          templateFile: 'tools/template/plop-go-templates/go-query/query.hbs',
          data: templateData
        }
      ];
    }
  });

  // --- Generator for a single Internal Repository (Implementation) ---
  plop.setGenerator('internal-repository', {
    description: 'Generates a new Go repository implementation',
    prompts: [moduleNamePrompt],
    actions: function(data) {
      return [
        {
          type: 'add',
          path: 'apps/users/internal/infra/persistence/repositories/{{kebabCase data.moduleName}}_repository.go',
          templateFile: 'tools/template/plop-go-templates/go-repository/repository.hbs',
          data: { moduleName: data.moduleName }
        }
      ];
    }
  });

  // --- Generator for Domain files ---
  plop.setGenerator('internal-domain', {
    description: 'Generates core domain files (Entity, Errors, Event, Repository Interface, ID)',
    prompts: [moduleNamePrompt],
    actions: function(data) {
      return [
        {
          type: 'add',
          path: 'apps/users/internal/domain/shared/ids/{{kebabCase data.moduleName}}_id.go',
          templateFile: 'tools/template/plop-go-templates/go-ids/id.hbs',
          data: { moduleName: data.moduleName }
        },
        {
          type: 'add',
          path: 'apps/users/internal/domain/{{kebabCase data.moduleName}}/{{kebabCase data.moduleName}}.go',
          templateFile: 'tools/template/plop-go-templates/go-domain/entity.hbs',
          data: { moduleName: data.moduleName }
        },
        {
          type: 'add',
          path: 'apps/users/internal/domain/{{kebabCase data.moduleName}}/repository.go',
          templateFile: 'tools/template/plop-go-templates/go-domain/repository_interface.hbs',
          data: { moduleName: data.moduleName }
        },
        {
          type: 'add',
          path: 'apps/users/internal/domain/{{kebabCase data.moduleName}}/errors.go',
          templateFile: 'tools/template/plop-go-templates/go-domain/errors.hbs',
          data: { moduleName: data.moduleName }
        },
        {
          type: 'add',
          path: 'apps/users/internal/domain/{{kebabCase data.moduleName}}/event.go',
          templateFile: 'tools/template/plop-go-templates/go-domain/event.hbs',
          data: { moduleName: data.moduleName }
        }
      ];
    }
  });

};