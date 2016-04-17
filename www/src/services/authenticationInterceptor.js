/**
 * authenticationInterceptor intercept every GET, PUSH, etc http PDUs
 * by adding the token user get when they log in (if they have any).
 * It will authorize users to get certain privileges.
 */
app.factory('authenticationInterceptor', ['$rootScope', '$q', '$window',
                                 function ($rootScope,   $q,   $window) {
  return {
    request: function (config) {
      config.headers = config.headers || {};
      if ($window.sessionStorage.token) {
        config.headers.Authorization = 'Bearer ' + $window.sessionStorage.token;
      }
      return config;
    },
    response: function (response) {
      if (response.status === 401) {
        // handle the case where the user is not authenticated
        console.log("Status code: 401");
      }
      return response || $q.when(response);
    }
  };
}]);

app.config(['$httpProvider', function ($httpProvider) {
  $httpProvider.interceptors.push('authenticationInterceptor');
}]);