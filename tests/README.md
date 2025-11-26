# Reasoning behind this package

Usually in Go we would like to keep tests close to the code they're testing.
However, in this case our tests have an external dependency (local docker container with local DynamoDb instance) that
needs to be setup & torn down for every test in the project,
which would require duplicating this setup/teardown logic for every package, which is error-prone.

Putting all tests in 1 package allows us to avoid that duplication and centralize the setup/teardown logic.

# Supposed structure for this package
    - package name
        - handler
        - repository
        - service
 