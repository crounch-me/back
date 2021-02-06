@account
Feature: Logout

  Scenario: OK
    Given I authenticate with a random user
    When I send a "POST" request on "/account/logout"
    Then the status code is 204
