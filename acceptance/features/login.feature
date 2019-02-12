Feature: Login

  Scenario: OK
    Given I create these user
      | email          | password |
      | login@test.com | test     |
    And I use this body
      """
        {
          "email": "login@test.com",
          "password": "test"
        }
      """ 
    When I send a "POST" request on "/users/login"
    Then the status code is 201
    And "$.accessToken" is a non empty string