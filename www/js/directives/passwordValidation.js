/* In development, so far just a copy paste */

app.directive('passwordValidation', function() {
   return {
      require: 'ngModel',
      scope: {
        passwordValidation: '='
      },
      link: function(scope, element, attrs, ctrl) {
        scope.$watch(function() {
            var combined;

            if (scope.passwordVerify || ctrl.$viewValue) {
               combined = scope.passwordVerify + '_' + ctrl.$viewValue; 
            }                    
            return combined;
        }, function(value) {
            if (value) {
                ctrl.$parsers.unshift(function(viewValue) {
                    var origin = scope.passwordVerify;
                    if (origin !== viewValue) {
                        ctrl.$setValidity('passwordValidation', false);
                        return undefined;
                    } else {
                        ctrl.$setValidity('passwordValidation', true);
                        return viewValue;
                    }
                });
            }
        });
     }
   };
});