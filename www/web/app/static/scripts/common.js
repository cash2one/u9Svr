requirejs.config({
  	baseUrl: '/static/scripts',
	paths: {
		json: 'lib/json2',
		mock: 'lib/mock',
		jquery: 'lib/jquery',
		mockjax: 'lib/jquery.mockjax',
		bootstrap:'lib/bootstrap',
		moment: 'lib/moment-with-locales',
		datepicker: 'lib/bootstrap-datetimepicker.min',
		paginator: 'lib/bootstrap-paginator.min',
		spin: 'lib/spin',
		//underscore: 'underscore',
		//backbone: 'backbone',
		//backboneLocalstorage: 'backbone.localStorage',
		//text: 'text',
	},
	shim: {
		'paginator': {
			deps: ['jquery']
		},
		'bootstrap':{
			deps:['jquery']
		}
	}
	// shim: {
	// 	underscore: {
	// 		exports: '_'
	// 	},
	// 	backbone: {
	// 		deps: [
	// 			'underscore',
	// 			'jquery'
	// 		],
	// 		exports: 'Backbone'
	// 	},
	// 	backboneLocalstorage: {
	// 		deps: ['backbone'],
	// 		exports: 'Store'
	// 	}
	// }
	// config: {
	//     moment: {noGlobal: true}
	// }
})