/**
 * ProfileEditController handles the edit page for the user profile.
 * Will get info from server, handle edits, and then push changes back to server
 */
app.controller('ProfileEditController', ['$scope', '$http', '$window', '$location', '$timeout', '$interval', 'tokenRefresher', 'toastr',
                                function ($scope,   $http,   $window,   $location,   $timeout,   $interval,   tokenRefresher,   toastr) {
  // Declare variables
  $scope.user = { // Placeholder
      FullName: 'Full Name',
      EMail: 'Your email adress',
      Phone: 'Your phone number',
      Description: 'Input your description here',

      ProfileIcon: 'img/profileDefault.png',
      ProfileHeader: 'img/backgroundDefault.png',

      PDFs: [
      ]
  };
  $scope.message = '';
  $scope.currentPDF = -1; // no pdfs in array
  $scope.saved = true;
  $scope.loading = {header: false, icon: false, pdf: false};
  $scope.maxLength = {
    fullName: 70,
    email: 50,
    phone: 80,
    description: 150,
    profileIcon: 150,
    profileHeader: 360,
    pdfs: 21844
  }
  $scope.logButton = {title: 'Log out', 
    click: function() {
      if ($window.sessionStorage.getItem('token') != null) {
        
        $http.get('/api/logout').then(
          function success(response) {
            $window.sessionStorage.removeItem('token');
            toastr.success('You have been logged out');
            $location.path('/');
          },
          function error(response) {
            $window.sessionStorage.removeItem('token');
            toastr.success('You have been logged out');
            $location.path('/');
          }
        );
      } else {
        $location.path('/login');
      }
    }
  };
  

  /* Notify user that changes not saved (on window close) */
  var checkIfSaved = function() {
    $scope.$watch('user', function(newValue, oldValue) {
        if (newValue != oldValue) {
          $scope.saved = false;
        }
      }, true); // true to check for value equality (not reference)
    
    var preventClosingIfNotSaved = function(event) {
      if ($scope.saved == false) {
        event.returnValue = "Warning. Leaving without publishing will remove changes."
      }
    };
    if (window.addEventListener) {
      window.addEventListener("beforeunload", preventClosingIfNotSaved);
    } else {
      // For IE browsers
      window.attachEvent("onbeforeunload", preventClosingIfNotSaved);
    }
  };

  /* Do a http request to server*/
  $http.get('/api/profile/get-edit').then(

    // Get user information from server and puts it in the user variable
    function success(response) {
      if (response.data != '') {
        $scope.user = response.data;
        if ($scope.user.PDFs.length > 0) {
          $scope.currentPDF = 0;
        }
        checkIfSaved(); // start checking for saves to prevent closing without saving
      }
    },
    function error(response) {
      if (response.status == 401) {
        $scope.message = 'You don\'t have permission to access this content';
        toastr.warning($scope.message+'. Please log in to edit profile');
        $location.path('/login');
        

      } else if (response.status == 400) {
        $scope.message = 'You are not logged in';
        toastr.warning($scope.message+'. Please log in to edit profile');
        $location.path('/'); // return to start page

      } else if (response.status == 413) {

      } else {
        $scope.message = 'User is not found';
      }
    }
  );

  /* Publish changes to server */
  $scope.publish = function() {
    if ($scope.logButton.click == 'login()') {
      $location.path('/login');
      
    } else {
      
      $http.post('api/profile/save', $scope.user).then(
        function success(response) {
          var oldMessage = $scope.message;
          $scope.message = 'Saved Success';
          $scope.saved = true;
          
          // Change back message after 5 seconds:
          $interval(function() {$scope.message = oldMessage;}, 5*1000);
        },
        function error(response) {
          var oldMessage = $scope.message;
          $scope.message = 'Saved failed';
          $interval(function() {$scope.message = oldMessage;}, 5*1000); // 5 sec
        }
      );
    }
  };
  
  /* To happen if needing to log in again */
  tokenRefresher.subscribeOnReloginEvent($scope, function() {
    $scope.logButton.title = 'Log in';
    $scope.message = 'You have been logged out. Please log in again';
  });
  
  
  $scope.changePDF = function(n) {
    $scope.currentPDF = n;
  };
  $scope.changeBackground = function(response) {
    $scope.user.ProfileHeader = response.data;
  };
  $scope.changeProfileIcon = function(response) {
    $scope.user.ProfileIcon = response.data //This is never run?
  };
  $scope.addPDF = function(response) {
    if ($scope.user.PDFs == null) {
      $scope.user.PDFs = [];
    }
    $scope.user.PDFs.push({Title: 'Unnamed', Path: response.data});
    $scope.currentPDF = $scope.user.PDFs.length-1;
  };

  $scope.startUploadHeader = function(response) {
    $scope.loading.header = true;
  };
  $scope.startUploadIcon = function(response) {
    $scope.loading.icon = true;
  };
  $scope.startUploadPdf = function(response) {
    $scope.loading.pdf = true;
  };
  $scope.endUploadHeader = function(response) {
    $scope.loading.header = false;
  };
  $scope.endUploadIcon = function(response) {
    $scope.loading.icon = false;
  };
  $scope.endUploadPDF = function(response) {
    $scope.loading.pdf = false;
  };
}]);
