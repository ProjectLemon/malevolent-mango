/**
 * ProfileController handles the public profile page. Gets info from server
 * to then be shown.
 */
app.controller('ProfileController', ['$scope', '$http', '$window', '$routeParams', function ($scope, $http, $window, $routeParams) {
  // Declare variables
  $scope.user = {
      FullName: 'Nathan Drake',
      Email: 'testing@example.com',
      Phone: '073-902301',
      Description: 'Some description about something or other. Oh look at me, I\'m just writing enough text to get a new line. Lorem ipsum dolar cofal(?) and all that shit.',

      ProfileIcon: 'img/testFace.png',
      ProfileHeader: 'img/testBG2.png',

      Pdfs: [
          {title: 'Portfolio', path: 'pdfs/portfolio1.pdf'},
          {title: 'Resumé', path: 'pdfs/resume1.pdf'},
          {title: $routeParams.userID}
      ]
  };
  $scope.message = '';


  //Do a http request to server
  $http.get('/api/profile/get').then(
    //If success
    // Get user information from server and puts it in the user variable
    function (response) {
      if (response.status == 204) {
        
      } else {
        if (response.data != '') {
          $scope.user = response.data;
        }
      }
    },
    //If Error
    function(response) {
      if (response.status == 401) {
        $scope.message = 'You don\'t have permission to access this content'; 
      } else {
        $scope.message = 'User is not found'; 
      }
    }
  )
}]);
