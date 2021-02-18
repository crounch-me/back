@account
Feature: Sign up

  Scenario: OK
    Given I use this body
      """
        {
          "email": "signup@test.com",
          "password": "test"
        }
      """
    When I send a "POST" request on "/account/signup"
    Then the status code is 204

  Scenario: KO - Invalid body
    Given I use an invalid body
    And I send a "POST" request on "/account/signup"
    Then the status code is 400
    And "$.error" has string value "unmarshal-error"

  Scenario: KO - Missing fields
    Given I use an empty valid body
    And I send a "POST" request on "/account/signup"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    # And "$.fields[0].name" has string value "email"
    # And "$.fields[0].error" has string value "required"
    # And "$.fields[1].name" has string value "password"
    # And "$.fields[1].error" has string value "required"

  Scenario: KO - Invalid fields
    Given I use this body
      """
        {
          "email": "a",
          "password": "a"
        }
      """
    And I send a "POST" request on "/account/signup"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    # And "$.fields[0].name" has string value "email"
    # And "$.fields[0].error" has string value "email"
    # And "$.fields[1].name" has string value "password"
    # And "$.fields[1].error" has string value "gt"

  Scenario: KO - Duplicated user
    Given I use this body
      """
        {
          "email": "duplicated@test.com",
          "password": "test"
        }
      """
    And I send a "POST" request on "/account/signup"
    And I use this body
      """
        {
          "email": "duPliCated@test.com",
          "password": "test"
        }
      """
    When I send a "POST" request on "/account/signup"
    Then the status code is 409
    And "$.error" has string value "duplicate-user-error"
