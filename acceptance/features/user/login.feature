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
