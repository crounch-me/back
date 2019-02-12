Feature: Sign up

  Scenario: OK
    Given I use this body
      """
        {
          "email": "signup@test.com",
          "password": "test"
        }
      """
    When I send a "POST" request on "/users"
    Then the status code is 201
    And "$.email" has string value "signup@test.com"

  Scenario: KO - Duplicated user
    Given I use this body
      """
        {
          "email": "duplicated@test.com",
          "password": "test"
        }
      """
    And I send a "POST" request on "/users"
    When I send a "POST" request on "/users"
    Then the status code is 409
    And "$.error" has string value "duplicate"
    And "$.errorDescription" has string value "User with this email already exists"