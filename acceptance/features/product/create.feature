Feature: Create Product

  Scenario: OK
    Given I authenticate with a random user
    And I create these products
      | name             |
      | Creation produit |
    Then the status code is 201
    And "$.id" is a non empty string
    And "$.name" is a string equal to "Creation produit"

  Scenario: KO - Invalid body
    Given I authenticate with a random user
    And I use an invalid body
    When I send a "POST" request on "/products"
    Then the status code is 400
    And "$.error" has string value "unmarshal-error"

  Scenario: KO - Missing fields
    Given I use an empty valid body
    And I send a "POST" request on "/products"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "name"
    And "$.fields[0].error" has string value "required"
