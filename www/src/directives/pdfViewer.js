/**
 * pdfViewer shows a pdf and update it when source change
 * Should have two attributes, pdfs and current. Current is only
 * a number but pdfs should be an array of pdfs in form of
 * [{Title: 'title', Path: 'som/path'}, ...]
 */
app.directive('pdfViewer', ['$compile', function($compile) {
  return {
    restrict: 'E',
    scope: {
      pdfs: '@',
      current: '='
    },
    
    link: function(scope, element, attrs) {
      
      var currentElement = element;
      scope.$parent.$watch(scope.pdfs, function(newValue, oldValue) {
        var html;
        if (newValue && newValue.length > 0) {
          var pdf = newValue[scope.current];
          if (pdf != null) {
            html = '<embed src="'+pdf.Path+'" type="application/pdf">'
            
          } else {
            html = '<p>No pdf</p>'
          } 
        } else {
          html = '<p>No pdf</p>'
        }
        var e = $compile(html)(scope);
        currentElement.replaceWith(e);
        currentElement = e;
      });
    }
  };
}]);