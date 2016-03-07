var app = angular.module('app', ['ngSanitize', 'ui.bootstrap', 'LocalStorageModule',]);

app.controller('HomeController',
  function($scope, $uibModal, localStorageService)
  {

    $scope.btnSaveClicked = function() {
      var modalInstance = $uibModal.open({
        animation: true,
        template:
          '<div class="modal-content">' +
          '  <div class="modal-header">' +
          '    <button type="button" class="close" ng-click="cancel()">x</button>' +
          '    <h4 class="modal-title" id="mySmallModalLabel">Save</h4>' +
          '  </div>' +
          '  <div class="modal-body">' +
          '    <input ng-model="input.name" type="text" class="form-control" ng-enter="btnSaveClicked()" focus-me required />' +
          '  </div>' +
          '  <div class="modal-footer">' +
          '    <input type="button" value="Save" class="btn btn-primary" ng-click="btnSaveClicked()" />' +
          '    <input type="button" value="Close" class="btn btn-default" ng-click="cancel()" />' +
          '  </div>' +
          '</div>',
        controller: function($scope, tree, $uibModalInstance) {
          $scope.input = {};
          $scope.btnSaveClicked = function() {
            var names = localStorageService.get('names') || {};
            names[$scope.input.name] = tree;
            localStorageService.set('names', names);
          };
          $scope.cancel = function() {
            $uibModalInstance.dismiss('cancel');
          };
        },
        size: 'sm',
        resolve: {
          tree: function() { return $scope.tree; }
        },
      });
    };

    $scope.btnLoadClicked = function() {
      var modalInstance = $uibModal.open({
        animation: true,
        template:
          '<div class="modal-content">' +
          '  <div class="modal-header">' +
          '    <button type="button" class="close" ng-click="cancel()">x</button>' +
          '    <h4 class="modal-title" id="mySmallModalLabel">Load</h4>' +
          '  </div>' +
          '  <div class="modal-body">' +
          '    <select ng-model="input.name" ng-options="name for name in names" class="form-control" ng-enter="btnLoadClicked()" focus-me></select>' +
          '  </div>' +
          '  <div class="modal-footer">' +
          '    <input type="button" value="Load" class="btn btn-primary" ng-click="btnLoadClicked()" />' +
          '    <input type="button" value="Close" class="btn btn-default" ng-click="cancel()" />' +
          '  </div>' +
          '</div>',
        controller: 'LoadModalController',
        size: 'sm',
      }).result.then(function(res) {
        if (res) {
          $scope.tree = res;
        }
      });
    };


    var highlightJson = function(json) {
      json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
      return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function(match) {
        var cls = 'number';
        if (/^"/.test(match)) {
          if (/:$/.test(match)) {
            cls = 'key';
          } else {
            cls = 'string';
          }
        } else if (/true|false/.test(match)) {
          cls = 'boolean';
        } else if (/null/.test(match)) {
          cls = 'null';
        }
        return '<span class="' + cls + '">' + match + '</span>';
      });
    };

    $scope.highlight = function(obj) {
      return highlightJson(angular.toJson(obj, true))
    };


    (function constructor() {
      $scope.selectedNode = null;
      $scope.tree = [];

      $scope.colors = [
        '#5aa1be',
        '#c065b8',
        '#c0ab5f',
        '#59bc95',
        '#7a63be',
        '#c0615b',
        '#73bb5e',
        '#6ee585',
        '#7088e8',
        '#eb77a3',
        '#c2e675',
      ];

    })();

  }
);
