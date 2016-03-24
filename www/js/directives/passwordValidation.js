/* In development, so far just a copy paste */
/*
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
*/
/*
app.directive('verifyPassword', function() {
    return {
        require: 'ngModel',
        link: function (scope, elem, attrs, model) {
            if (!attrs.verify) {
                console.error('verifyPassword expects a model as an argument!');
                return;
            }
            scope.$watch(attrs.verify, function (value) {
                // Only compare values if the second ctrl has a value.
                if (model.$viewValue !== undefined && model.$viewValue !== '') {
                    model.$setValidity('verifyPassword', value === model.$viewValue);
                }
            });
            model.$parsers.push(function (value) {
                // Mute the verifyPassword error if the second ctrl is empty.
                if (value === undefined || value === '') {
                    model.$setValidity('verifyPassword', true);
                    return value;
                }
                var isValid = value === scope.$eval(attrs.verify);
                model.$setValidity('verifyPassword', isValid);
                return isValid ? value : undefined;
            });
        }
    };
});
*/
/*
app.directive('sameAs', function() {
  return {
    require: 'ngModel',
    link: function(scope, elm, attrs, ctrl) {
      ctrl.$parsers.unshift(function(viewValue) {
        if (viewValue === scope[attrs.sameAs]) {
          ctrl.$setValidity('sameAs', true);
          return viewValue;
        } else {
          ctrl.$setValidity('sameAs', false);
          return undefined;
        }
      });
    }
  };
});
*/
app.directive("compareTo", function() {
    return {
        require: "ngModel",
        scope: {
            otherModelValue: "=compareTo"
        },
        link: function(scope, element, attributes, ngModel) {
             
            ngModel.$validators.compareTo = function(modelValue) {
                return modelValue == scope.otherModelValue;
            };
 
            scope.$watch("otherModelValue", function() {
                ngModel.$validate();
            });
        }
    };
})

