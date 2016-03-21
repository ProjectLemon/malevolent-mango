
app.controller('SignupController', ['$scope', '$http', '$window', function ($scope, $http, $window) {
  // Declare variables
  $scope.user = {};
  $scope.message = '';
  
  // Submit function
  $scope.submit = function () {
    /* If form is filled out, send data to server */
    if ($scope.user.email && $scope.user.password && 
        $scope.user.email != '' && $scope.user.password != '') {
  
      $http
        .post('/api/register', $scope.user)
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
            if (response.data.Token !=== undefined) {
              $window.sessionStorage.token = response.data.Token;
              $scope.message = 'Logged in';
              
            } else {
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