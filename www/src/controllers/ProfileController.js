/**
 * ProfileController handles the public profile page. Gets info from server
 * to then be shown.
 */
app.controller('ProfileController', ['$scope', '$http', '$window', '$location', '$routeParams', 'toastr',
                            function ($scope,   $http,   $window,   $location,   $routeParams,   toastr) {
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
  $scope.currentPDF = -1;
  $scope.message = '';


  //Do a http request to server
  $http.get('/api/profile/get-view/'+$routeParams.publicName).then(
    // If success
    // Get user information from server and puts it in the user variable
    function (response) {
      if (response.status == 204) { // No Content
        $scope.message = 'User '+$routeParams.publicName+' is not found';
        toastr.warning($scope.message);
        $location.path('/');
        
      } else if (response.data != '') {
        $scope.user = response.data;
        $scope.currentPDF = 0;
      }
    },
    // If Error
    function(response) {
      if (response.status == 401) {
        $scope.message = 'You don\'t have permission to access this content'; 
      } else {
        $scope.message = 'User is not found'; 
      }
      toastr.warning($scope.message);
      $location.path('/');
    }
  );
  
  $scope.changePDF = function(n) {
    $scope.currentPDF = n;
  };
}]);
