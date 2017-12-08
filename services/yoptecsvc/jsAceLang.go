package yoptecsvc

const jsAceLang=`ace.define("ace/mode/yoptec_highlight_rules",["require","exports","module","ace/lib/oop","ace/mode/text_highlight_rules"], function(require, exports, module) {
	"use strict";
	
	var oop = require("../lib/oop");
	var TextHighlightRules = require("./text_highlight_rules").TextHighlightRules;
	
	var yoptecHighlightRules = function() {
	
		var keywords = (
			"харэ|жы|захуярить|иличовжопураз|иличовжопуразвилкойвглаз|го|йопта|вилкойвглаз|чоунастут|по|естьчо|аеслинайду|апохуй|ассо|пацансделал|"+
			 "отвечаю|атоэто|пероподребро|потрещим|иличо|И|НЕ|хапнуть|пнх|гоп|двигай|петух|клеенка|"+
			 "стопэжы|Конецвилкойвглаз|стопэйопта|стопэхапать|Конецестьчоа"
		);
	
		var builtinConstants = ("чотко|нечотко|порожняк|NULL|ДлительностьНаносекунды|"+
			"ДлительностьМикросекунды|ДлительностьМиллисекунды|ДлительностьСекунды|"+
			"ДлительностьМинуты|ДлительностьЧаса|ДлительностьДня|АргументыЗапуска");
	
		var functions = (
			"Число|Строка|Булево|ЦелоеЧисло|Массив|Структура|Дата|Длительность|"+
			"Импорт|Длина|Диапазон|ТекущаяДата|ПрошлоВремениС|Пауза|Хэш|"+
			"УникальныйИдентификатор|ПолучитьМассивчоунастутПула|ВернутьМассивВПул|СлучайнаяСтрока|НРег|ВРег|"+
			"Формат|КодСимвола|ТипЗнч|Сообщить|СообщитьФ|ОбработатьГорутины|ЗагрузитьИВыполнить|"+
			"ОписаниеОшибки|ПеременнаяОкружения|СтрСодержит|СтрСодержитЛюбой|СтрКоличество|СтрНайти|"+
			"СтрНайтиЛюбой|СтрНайтиПоследний|СтрЗаменить|Окр"
		);
	
		var builtinTypes = ("ГруппаОжидания|Сервер|Клиент|ФайловаяБазаДанных");
	
		var keywordMapper = this.createKeywordMapper({
			"keyword": keywords,
			"support.function": functions,
			"support.type": builtinTypes,
			"constant.language": builtinConstants,
			"variable.language": "self"
		}, "identifier");
		
		this.$rules = {
			"start" : [
			{
				token : "comment",
				regex : "\\/\\/.*$"
			},
			{
				token : "comment",
				regex : "\\#.*$"
			},
			{
				token : "string", // single line
				regex : /"(?:[^"\\]|\\.)*?"/
			}, {
				token : "string", // raw
				regex : '`+"`"+`',
				next : "bqstring"
			}, {
				token : "constant.numeric", // hex
				regex : "0[xX][0-9a-fA-F]+\\b" 
			}, {
				token : "constant.numeric", // float
				regex : "[+-]?\\d+(?:(?:\\.\\d*)?(?:[eE][+-]?\\d+)?)?\\b"
			}, {
				token : ["keyword", "text", "entity.name.function"],
				regex : "(йопта)(\\s+)([a-zA-Zа-яА-ЯёЁ_$][a-zA-Zа-яА-Я0-9_$]*)(?![a-zA-Zа-яА-ЯёЁ])"
			}, {
				token : keywordMapper,
				regex : "[a-zA-Zа-яА-Я_$][a-zA-Zа-яА-Я0-9_$]*(?![a-zA-Zа-яА-ЯёЁ])"
			}, {
				token : "keyword.operator",
				regex : "!|\\$|%|&|\\*|\\-\\-|\\-|\\+\\+|\\+|~|==|=|!=|<=|>=|<<=|>>=|>>>=|<>|<|>|!|&&|\\|\\||\\:|\\*=|%=|\\+=|\\-=|&=|\\^="
			}, {
				token : "paren.lparen",
				regex : "[\\[\\(\\{]"
			}, {
				token : "paren.rparen",
				regex : "[\\]\\)\\}]"
			}, {
				token : "text",
				regex : "\\s+|\\w+"
			} ],
			"bqstring" : [
                {
                    token : "string",
                    regex : '`+"`"+`',
                    next : "start"
                }, {
                    defaultToken : "string"
                }
            ]
		};
		
		this.normalizeRules();
	}
	
	oop.inherits(yoptecHighlightRules, TextHighlightRules);
	
	exports.yoptecHighlightRules = yoptecHighlightRules;
	});
		
	ace.define("ace/mode/yoptec",["require","exports","module","ace/lib/oop","ace/mode/text","ace/mode/yoptec_highlight_rules","ace/mode/folding/yoptec","ace/range","ace/worker/worker_client"], function(require, exports, module) {
	"use strict";
	
	var oop = require("../lib/oop");
	var TextMode = require("./text").Mode;
	var yoptecHighlightRules = require("./yoptec_highlight_rules").yoptecHighlightRules;
	var Range = require("../range").Range;
	var WorkerClient = require("../worker/worker_client").WorkerClient;
	
	var Mode = function() {
		this.HighlightRules = yoptecHighlightRules;
		
		this.$behaviour = this.$defaultBehaviour;
	};
	oop.inherits(Mode, TextMode);
	
	(function() {
	   
		this.lineCommentStart = "//";
		
		var indentKeywords = {
			"йопта": 1,
			"атоэто": 1,
			"жы": 1,
			"иличовжопураз": 1,
			"иличовжопуразвилкойвглаз": 1,
			"аеслинайду": 1,
			"апохуй": 1,
			"потрещим": 1,
			"хапнуть": 1,
			"гоп": 1,
			"стопэжы": -1,
			"Конецвилкойвглаз": -1,
			"стопэйопта": -1,
			"стопэхапать": -1,
			"Конецестьчоа": -1
		};
		var outdentKeywords = [
			"иличовжопураз",
			"иличовжопуразвилкойвглаз",
			"стопэжы",
			"Конецвилкойвглаз",
			"стопэйопта",
			"стопэхапать",
			"Конецестьчоа"
		];
	
		function getNetIndentLevel(tokens) {
			var level = 0;
			for (var i = 0; i < tokens.length; i++) {
				var token = tokens[i];
				if (token.type == "keyword") {
					if (token.value in indentKeywords) {
						level += indentKeywords[token.value];
					}
				} else if (token.type == "paren.lparen") {
					level += token.value.length;
				} else if (token.type == "paren.rparen") {
					level -= token.value.length;
				}
			}
			if (level < 0) {
				return -1;
			} else if (level > 0) {
				return 1;
			} else {
				return 0;
			}
		}
	
		this.getNextLineIndent = function(state, line, tab) {
			var indent = this.$getIndent(line);
			var level = 0;
	
			var tokenizedLine = this.getTokenizer().getLineTokens(line, state);
			var tokens = tokenizedLine.tokens;
	
			if (state == "start") {
				level = getNetIndentLevel(tokens);
			}
			if (level > 0) {
				return indent + tab;
			} else if (level < 0 && indent.substr(indent.length - tab.length) == tab) {
				if (!this.checkOutdent(state, line, "\n")) {
					return indent.substr(0, indent.length - tab.length);
				}
			}
			return indent;
		};
	
		this.checkOutdent = function(state, line, input) {
			if (input != "\n" && input != "\r" && input != "\r\n")
				return false;
	
			if (line.match(/^\s*[\)\}\]]$/))
				return true;
	
			var tokens = this.getTokenizer().getLineTokens(line.trim(), state).tokens;
	
			if (!tokens || !tokens.length)
				return false;
	
			return (tokens[0].type == "keyword" && outdentKeywords.indexOf(tokens[0].value) != -1);
		};
	
		this.autoOutdent = function(state, session, row) {
			var prevLine = session.getLine(row - 1);
			var prevIndent = this.$getIndent(prevLine).length;
			var prevTokens = this.getTokenizer().getLineTokens(prevLine, "start").tokens;
			var tabLength = session.getTabString().length;
			var expectedIndent = prevIndent + tabLength * getNetIndentLevel(prevTokens);
			var curIndent = this.$getIndent(session.getLine(row)).length;
			if (curIndent <= expectedIndent) {
				return;
			}
			session.outdentRows(new Range(row, 0, row + 2, 0));
		};
	
		this.createWorker = function(session) {
			var worker = new WorkerClient(["ace"], "ace/mode/yoptec_worker", "Worker");
			worker.attachToDocument(session.getDocument());
			
			worker.on("annotate", function(e) {
				session.setAnnotations(e.data);
			});
			
			worker.on("terminate", function() {
				session.clearAnnotations();
			});
			
			return worker;
		};
	
		this.$id = "ace/mode/yoptec";
	}).call(Mode.prototype);
	
	exports.Mode = Mode;
	});
`