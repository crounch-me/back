@product-in-list
Feature: Update a product in a list

  Scenario: OK
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I create these products
      | name                |
      | Mon premier produit |
    And I send a "POST" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}"
    When I send a "PATCH" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}/buy"
    Then the status code is 204

  Scenario: KO - List id not uuid
    Given I authenticate with a random user
    When I send a "PATCH" request on "/listing/lists/a/products/{{ .ProductID }}/buy"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "listID"
    And "$.fields[0].error" has string value "uuid"

  Scenario: KO - Product id not uuid
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    When I send a "PATCH" request on "/listing/lists/{{ .ListID }}/products/a/buy"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "productID"
    And "$.fields[0].error" has string value "uuid"

  Scenario: KO - User not authorized
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I create these products
      | name                |
      | Mon premier produit |
    And I authenticate with a random user
    When I send a "PATCH" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}/buy"
    Then the status code is 403
    And "$.error" has string value "forbidden-error"

  Scenario: KO - Product not in list
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I create these products
      | name                |
      | Mon premier produit |
    When I send a "PATCH" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}/buy"
    Then the status code is 404
    And "$.error" has string value "product-not-in-list-error"

