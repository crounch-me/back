Feature: Login

  Scenario: OK
    Given I authenticate with a random user
    When I send a "POST" request on "/logout"
    And I send a "GET" request on "/me"
    Then the status code is 401
    And "$.error" is a string equal to "unauthorized-error"
