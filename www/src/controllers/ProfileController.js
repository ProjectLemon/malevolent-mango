/**
 * ProfileController handles the public profile page. Gets info from server
 * to then be shown.
 */
app.controller('ProfileController', ['$scope', '$http', '$window', '$routeParams', 
                            function ($scope,   $http,   $window,   $routeParams) {
  // Declare variables
  $scope.user = {
      FullName: '',
      EMail: '',
      Phone: '',
      Description: '',

      ProfileIcon: '',
      ProfileHeader: '',

      PDFs: [
      ]
  };
  $scope.message = '';


  //Do a http request to server
  $http.post('/api/profile/get-view', {UserID: $routeParams.userID}).then(
    // If success
    // Get user information from server and puts it in the user variable
    function (response) {
      if (response.status == 204) {
        $scope.message = 'User is not found';
        $window.location.href = '#/'
        
      } else if (response.data != '') {
        $scope.user = response.data;
      }
    },
    // If Error
    function(response) {
      if (response.status == 401) {
        $scope.message = 'You don\'t have permission to access this content'; 
      } else {
        $scope.message = 'User is not found'; 
      }
    }
  )
}]);
