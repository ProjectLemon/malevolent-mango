app.controller('ProfileController', ['$scope', '$http', '$window', function ($scope, $http, $window) {
    // Declare variables
    $scope.user = {};
    $scope.message = ''; 

    
    //Do a http request to server
    $http.get('/profile').then(
        //If success
        // Get user information from server and puts it in the user variable
        function (response) {
            $scope.user = response.data;
        },
        //If Error
        // Display message that the user is not found
        function(response){
            $scope.message = 'User is not found'; 
        }
    )
    }]);