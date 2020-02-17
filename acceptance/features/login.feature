Feature: Login

  Scenario: CORS
    When I send a "OPTIONS" request on "/users/login"
    Then the status code is 200
    And the header "Access-Control-Allow-Methods" equals "POST,OPTIONS"
    And the header "Access-Control-Allow-Origin" equals "*"
    And the header "Access-Control-Allow-Headers" equals "Content-Type,Authorization,Accept"

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
    Then the status code is 403
