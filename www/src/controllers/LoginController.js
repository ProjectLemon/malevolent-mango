/**
 * LoginController handles the login form, and send the data to server for confirmation
 */
app.controller('LoginController', ['$scope', '$http', '$window', function ($scope, $http, $window) {
  // Declare variables
  $scope.user = {};
  $scope.message = '';
  
  /* Submit function, sends login info (email & password) to server */
  $scope.submit = function () {
  
    /* If form is filled out, send data to server */
    if ($scope.user.email && $scope.user.password && 
        $scope.user.email != '' && $scope.user.password != '') {
  
      $http
        .post('/api/login', $scope.user)
        .then(
          /* On success */
          function (response) {
            /* The response object has these properties:
  
            • data – {string|Object} – The response body transformed with the transform functions.
            • status – {number} – HTTP status code of the response.
            • headers – {function([headerName])} – Header getter function.
            • config – {Object} – The configuration object that was used to generate the request.
            • statusText – {string} – HTTP status text of the response.
            */
            if (response.data.Token !== undefined) {
              /* Complete success */
              $window.sessionStorage.token = response.data.Token;
              $scope.message = 'Logged in';
              $window.location.href = '/';
              
            } else {
              /* No server response */
              $scope.message = 'Unable to contact server';
            }
          },
          /* On error */
          function (response) {
            // Erase the token if the user fails to log in
            delete $window.sessionStorage.token;
  
            // Handle login errors here
            $scope.message = response.data;
          }
        );
    } else {
      $scope.formNotFilled = true; // to signify to form
      $scope.message = 'Please fill out the form correctly';
    }
  };
}]);