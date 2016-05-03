/**
 * Redirects users trying to navigate to a unused url: #/profile
 */
app.controller('ProfileRedirectController', ['$window', '$location',
                                    function ($window,   $location) {
  if ($window.sessionStorage.getItem('token')) {
    $location.path('/profile/edit');
  } else {
    $location.path('/');
  }
}]);