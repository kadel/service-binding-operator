@naming-strategy-all
Feature: Bind an application to a service using custom naming strategies

    As a user of Service Binding Operator
    I want to bind applications to services it depends on

    Background:
        Given Namespace [TEST_NAMESPACE] is used
        * Service Binding Operator is running
        * PostgreSQL DB operator is installed

    Scenario: Bind an imported Node.js application to PostgreSQL database in the following order: Application, DB and Service Binding and Naming Strategy
        Given Imported Nodejs application "nodejs-rest-http-crud-naming-strategy" is running
        * DB "db-demo-naming-strategy" is running
        When Service Binding is applied
            """
            apiVersion: operators.coreos.com/v1alpha1
            kind: ServiceBinding
            metadata:
                name: binding-request-naming-strategy
            spec:
                application:
                    name: nodejs-rest-http-crud-naming-strategy
                    group: apps
                    version: v1
                    resource: deployments
                services:
                -   group: postgresql.baiju.dev
                    version: v1alpha1
                    kind: Database
                    name: db-demo-naming-strategy
            """
        Then Service Binding "binding-request-naming-strategy" is ready
        And application should be re-deployed
        And Secret "binding-request-naming-strategy" contains "DATABASE_DBNAME" key with value "db-demo-naming-strategy"
        And Secret "binding-request-naming-strategy" contains "DATABASE_USER" key with value "postgres"
        And Secret "binding-request-naming-strategy" contains "DATABASE_PASSWORD" key with value "password"
        And Secret "binding-request-naming-strategy" contains "DATABASE_DB_PASSWORD" key with value "password"
        And Secret "binding-request-naming-strategy" contains "DATABASE_DB_NAME" key with value "db-demo-naming-strategy"
        And Secret "binding-request-naming-strategy" contains "DATABASE_DB_PORT" key with value "5432"
        And Secret "binding-request-naming-strategy" contains "DATABASE_DB_USER" key with value "postgres"
        And Secret "binding-request-naming-strategy" contains "DATABASE_DB_HOST" key with dynamic IP addess as the value
        And Secret "binding-request-naming-strategy" contains "DATABASE_DBCONNECTIONIP" key with dynamic IP addess as the value
        And Secret "binding-request-naming-strategy" contains "DATABASE_DBCONNECTIONPORT" key with value "5432"


    @strategy-none
    Scenario: Bind an imported Node.js application to PostgreSQL database in the following order: Application, DB and Service Binding and Naming Strategy none
        Given Imported Nodejs application "nodejs-rest-http-crud-naming-none" is running
        * DB "db-demo-naming-none" is running
        When Service Binding is applied
            """
            apiVersion: operators.coreos.com/v1alpha1
            kind: ServiceBinding
            metadata:
                name: binding-request-naming-none
            spec:
                namingStrategy: none
                application:
                    name: nodejs-rest-http-crud-naming-none
                    group: apps
                    version: v1
                    resource: deployments
                services:
                -   group: postgresql.baiju.dev
                    version: v1alpha1
                    kind: Database
                    name: db-demo-naming-none
            """
        Then Service Binding "binding-request-naming-none" is ready
        And application should be re-deployed
        And Secret "binding-request-naming-none" contains "dbname" key with value "db-demo-naming-none"
        And Secret "binding-request-naming-none" contains "user" key with value "postgres"
        And Secret "binding-request-naming-none" contains "password" key with value "password"
        And Secret "binding-request-naming-none" contains "db_password" key with value "password"
        And Secret "binding-request-naming-none" contains "db_name" key with value "db-demo-naming-none"
        And Secret "binding-request-naming-none" contains "db_port" key with value "5432"
        And Secret "binding-request-naming-none" contains "db_user" key with value "postgres"
        And Secret "binding-request-naming-none" contains "db_host" key with dynamic IP addess as the value
        And Secret "binding-request-naming-none" contains "dbConnectionIP" key with dynamic IP addess as the value
        And Secret "binding-request-naming-none" contains "dbConnectionIP" key with value "5432"


    Scenario: Bind an imported Node.js application to PostgreSQL database in the following order: Application, DB and Service Binding and custom Naming Strategy
        Given Imported Nodejs application "nodejs-rest-http-crud-naming-custom" is running
        * DB "db-demo-naming-custom" is running
        When Service Binding is applied
            """
            apiVersion: operators.coreos.com/v1alpha1
            kind: ServiceBinding
            metadata:
                name: binding-request-naming-custom
            spec:
                namingStrategy: "DB_{{ .name | upper }}_ENV"
                application:
                    name: nodejs-rest-http-crud-naming-custom
                    group: apps
                    version: v1
                    resource: deployments
                services:
                -   group: postgresql.baiju.dev
                    version: v1alpha1
                    kind: Database
                    name: db-demo-naming-custom
            """
        Then Service Binding "binding-request-naming-custom" is ready
        And application should be re-deployed
        And Secret "binding-request-naming-custom" contains "DB_DBNAME_ENV" key with value "db-demo-naming-custom"
        And Secret "binding-request-naming-custom" contains "DB_USER_ENV" key with value "postgres"
        And Secret "binding-request-naming-custom" contains "DB_PASSWORD_ENV" key with value "password"
        And Secret "binding-request-naming-custom" contains "DB_DB_PASSWORD_ENV" key with value "password"
        And Secret "binding-request-naming-custom" contains "DB_DB_NAME_ENV" key with value "db-demo-naming-custom"
        And Secret "binding-request-naming-custom" contains "DB_DB_PORT_ENV" key with value "5432"
        And Secret "binding-request-naming-custom" contains "DB_DB_USER_ENV" key with value "postgres"
        And Secret "binding-request-naming-custom" contains "DB_DB_HOST_ENV" key with dynamic IP addess as the value
        And Secret "binding-request-naming-custom" contains "DB_DBCONNECTIONIP_ENV" key with dynamic IP addess as the value
        And Secret "binding-request-naming-custom" contains "DB_DBCONNECTIONPORT_ENV" key with value "5432"

    Scenario: Bind an imported Node.js application to PostgreSQL database in the following order: Application, DB and Service Binding and Naming Strategy for bind as a file
        Given Imported Nodejs application "nodejs-rest-http-crud-naming-files" is running
        * DB "db-demo-naming-files" is running
        When Service Binding is applied
            """
            apiVersion: operators.coreos.com/v1alpha1
            kind: ServiceBinding
            metadata:
                name: binding-request-naming-files
            spec:
                bindAsFiles: true
                application:
                    name: nodejs-rest-http-crud-naming-files
                    group: apps
                    version: v1
                    resource: deployments
                services:
                -   group: postgresql.baiju.dev
                    version: v1alpha1
                    kind: Database
                    name: db-demo-naming-files
            """
        Then Service Binding "binding-request-naming-files" is ready
        And application should be re-deployed
        And Secret "binding-request-naming-files" contains "dbName" key with value "db-demo-naming-files"
        And Secret "binding-request-naming-files" contains "user" key with value "postgres"
        And Secret "binding-request-naming-files" contains "password" key with value "password"
        And Secret "binding-request-naming-files" contains "db_password" key with value "password"
        And Secret "binding-request-naming-files" contains "db_name" key with value "db-demo-naming-files"
        And Secret "binding-request-naming-files" contains "db_port" key with value "5432"
        And Secret "binding-request-naming-files" contains "db_user" key with value "postgres"
        And Secret "binding-request-naming-files" contains "db_host" key with dynamic IP addess as the value
        And Secret "binding-request-naming-files" contains "dbConnectionIP" key with dynamic IP addess as the value
        And Secret "binding-request-naming-files" contains "dbConnectionPort" key with value "5432"
