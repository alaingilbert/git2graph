var app = angular.module('app');

app.controller('LoadModalController',
  function($scope, $uibModalInstance, localStorageService)
  {

    $scope.cancel = function() {
      $uibModalInstance.dismiss('cancel');
    };


    $scope.btnLoadClicked = function() {
      $uibModalInstance.close($scope.projects[$scope.input.name]);
    };


    (function constructor() {
      $scope.projects = localStorageService.get('names') || {};
      $scope.names = _.chain($scope.projects).keys().sortBy().value();
      $scope.input = {};
      $scope.input.name = $scope.names[0];
    })();

  }
);
