
app.controller('LoginController', ['$scope', '$http', '$window', function ($scope, $http, $window) {
  $scope.user = {email: '', password: ''};
  $scope.message = '';
  $scope.submit = function () {
    console.log('Posting');
    $http
      .post('/login', $scope.user)
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
        
          $window.sessionStorage.token = response.data.token;
          $scope.message = 'Welcome';
        },
        /* On error */
        function (response) {
          // Erase the token if the user fails to log in
          delete $window.sessionStorage.token;

          // Handle login errors here
          $scope.message = response.data;
        }
      );
  };
}]);