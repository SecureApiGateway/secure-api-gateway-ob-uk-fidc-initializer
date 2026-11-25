/*
 * Script Type: OAuth2 Dynamic Client Registration
 * Purpose:
 */
(function () {
  if (operation === "CREATE" || operation === "UPDATE") {
    // Optional: Handle property the operation
    if (softwareStatement != null) {
        // Extract the jwks_uri from the SSA claims
        // Note: The ssa object allows access to decoded claims
        var softwareJwksEndpoint = softwareStatement.get("software_jwks_endpoint")
        if (softwareJwksEndpoint != null) {
            // 3. Store it in the clientProperties map
            // This map is what AM uses to save the client configuration to the CTS/Repo
            clientIdentity.addAttribute("oauth2ClientCustomProperties", "jwks_uri=" + softwareJwksEndpoint)
            clientIdentity.store()
        } else {
            // Optionally: Throw an exception to reject registration if this field is mandatory
            logger.warn("DCR Script: jwks_uri is missing from the provided SSA.")
        }
    } else {
        // Optional: throw an exception to prevent invalid token generation
        logger.warn("DCR Script: No software statement found in registration request.")
    }
  }
    // No return value is expected. Let it be undefined.
}());