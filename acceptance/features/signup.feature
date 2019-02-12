Feature: Sign up

  Scenario: OK
    Given I use this body
      """
        {
          "email": "test@test.com",
          "password": "test"
        }
      """
    When I send a "POST" request on "/users"
    Then "$.email" has string value "test@test.com"