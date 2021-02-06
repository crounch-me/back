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
    And I use this body
      """
        {
          "productId": "{{ .ProductID }}",
          "listId": "{{ .ListID }}"
        }
      """
    And I send a "POST" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}"
    And I use this body
      """
        {
          "productID": "{{ .ProductID }}",
          "listID": "{{ .ListID }}"
        }
      """
    When I send a "PATCH" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}/buy"
    Then the status code is 200
    And "$.listId" has string value "{{ .ListID }}"
    And "$.productId" has string value "{{ .ProductID }}"
    And "$.bought" has bool value "true"
    # And I send a "GET" request on "/lists/{{ .ListID }}"
    # And the returned products from list are
    #   | ID               | Name                | Category name | Bought |
    #   | {{ .ProductID }} | Mon premier produit | Divers        | Yes    |

  Scenario: KO - Invalid body
    Given I authenticate with a random user
    And I use an invalid body
    When I send a "PATCH" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}/buy"
    Then the status code is 400
    And "$.error" has string value "unmarshal-error"

  Scenario: KO - List not found
    Given I authenticate with a random user
    And I use this body
      """
        {
          "listID": "00000000-0000-0000-0000-000000000000",
          "productID": "{{ .ProductID }}"
        }
      """
    When I send a "PATCH" request on "/listing//lists/{{ .ListID }}/products/{{ .ProductID }}/buy"
    Then the status code is 404
    And "$.error" has string value "list-not-found-error"

  Scenario: KO - Product not found
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I use this body
      """
        {
          "listID": "{{ .ListID }}",
          "productID": "00000000-0000-0000-0000-000000000000"
        }
      """
    When I send a "PATCH" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}/buy"
    Then the status code is 404
    And "$.error" has string value "product-not-found-error"

  Scenario: KO - List id not uuid
    Given I authenticate with a random user
    And I use this body
      """
        {
          "listID": "a",
          "productID": "00000000-0000-0000-0000-000000000000"
        }
      """
    When I send a "PATCH" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}/buy"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "listID"
    And "$.fields[0].error" has string value "uuid"

  Scenario: KO - Product id not uuid
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    And I use this body
      """
        {
          "listID": "{{ .ListID }}",
          "productID": "a"
        }
      """
    When I send a "PATCH" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}/buy"
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
    And I use this body
      """
        {
          "listID": "{{ .ListID }}",
          "productID": "{{ .ProductID }}"
        }
      """
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
    And I use this body
      """
        {
          "listID": "{{ .ListID }}",
          "productID": "{{ .ProductID }}"
        }
      """
    When I send a "PATCH" request on "/listing/lists/{{ .ListID }}/products/{{ .ProductID }}/buy"
    Then the status code is 404
    And "$.error" has string value "product-in-list-not-found"

