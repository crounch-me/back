Feature: Login

  Scenario: OK
    Given I create these user
      | email          | password |
      | login@test.com | test     |
    And I use this body
      """
        {
          "email": "LoGiN@test.com",
          "password": "test"
        }
      """
    When I send a "POST" request on "/users/login"
    Then the status code is 201
    And "$.accessToken" is a non empty string

  Scenario: KO - Invalid body
    Given I use an invalid body
    And I send a "POST" request on "/users/login"
    Then the status code is 400
    And "$.error" has string value "unmarshal-error"

  Scenario: KO - Missing fields
    Given I use an empty valid body
    And I send a "POST" request on "/users/login"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "email"
    And "$.fields[0].error" has string value "required"
    And "$.fields[1].name" has string value "password"
    And "$.fields[1].error" has string value "required"

  Scenario: KO - Invalid fields
    Given I use this body
      """
        {
          "email": "a",
          "password": "a"
        }
      """
    And I send a "POST" request on "/users/login"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "email"
    And "$.fields[0].error" has string value "email"
    And "$.fields[1].name" has string value "password"
    And "$.fields[1].error" has string value "gt"

  Scenario: KO - User not found
    Given I use this body
      """
        {
          "email": "unknown-guy@test.com",
          "password": "test"
        }
      """
    When I send a "POST" request on "/users/login"
    Then the status code is 404
