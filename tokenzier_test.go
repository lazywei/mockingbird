package mockingbird_test

import (
	"io/ioutil"

	. "github.com/lazywei/mockingbird"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tokenzier", func() {
	Describe("extract tokens", func() {

		It("should skip string literals", func() {
			expectedResult := []string{"print"}

			Expect(ExtractTokens(`print ""`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print "Josh"`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print 'Josh'`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print "Hello \"Josh\""`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print 'Hello \'Josh\''`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print "Hello", "Josh"`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print 'Hello', 'Josh'`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print "Hello", "", "Josh"`)).To(Equal(expectedResult))
			Expect(ExtractTokens(`print 'Hello', '', 'Josh'`)).To(Equal(expectedResult))
		})

		It("should skip number literals", func() {
			Expect(ExtractTokens(`1 + 1`)).To(Equal([]string{`+`}))
			Expect(ExtractTokens(`add(123, 45)`)).To(Equal([]string{`add`, `(`, `)`}))
			Expect(ExtractTokens(`0x01 | 0x10`)).To(Equal([]string{`|`}))
			Expect(ExtractTokens(`500.42 * 1.0`)).To(Equal([]string{`*`}))
		})

		It("should extract common operators", func() {
			Expect(ExtractTokens("1 + 1")).To(Equal([]string{`+`}))
			Expect(ExtractTokens("1 - 1")).To(Equal([]string{`-`}))
			Expect(ExtractTokens("1 * 1")).To(Equal([]string{`*`}))
			Expect(ExtractTokens("1 / 1")).To(Equal([]string{`/`}))
			Expect(ExtractTokens("2 % 5")).To(Equal([]string{`%`}))
			Expect(ExtractTokens("1 & 1")).To(Equal([]string{`&`}))
			Expect(ExtractTokens("1 && 1")).To(Equal([]string{`&&`}))
			Expect(ExtractTokens("1 | 1")).To(Equal([]string{`|`}))
			Expect(ExtractTokens("1 || 1")).To(Equal([]string{`||`}))
			Expect(ExtractTokens("1 < 0x01")).To(Equal([]string{`<`}))
			Expect(ExtractTokens("1 << 0x01")).To(Equal([]string{`<<`}))
		})

		It("should skip comments", func() {

			cmtTests := []struct {
				str    string
				tokens []string
			}{
				{"foo\n# Comment", []string{`foo`}},
				{"foo\n# Comment\nbar", []string{`foo`, `bar`}},
				{"foo\n// Comment", []string{`foo`}},
				{"foo\n-- Comment", []string{`foo`}},
				{"foo\n\" Comment", []string{`foo`}},
				{"foo /* Comment */", []string{`foo`}},
				{"foo /* \nComment\n */", []string{`foo`}},
				{"foo <!-- Comment -->", []string{`foo`}},
				{"foo {- Comment -}", []string{`foo`}},
				{"foo (* Comment *)", []string{`foo`}},
				{"2 % 10\n% Comment", []string{`%`}},
				{"foo\n\"\"\"\nComment\n\"\"\"\nbar", []string{`foo`, `bar`}},
				{"foo\n'''\nComment\n'''\nbar", []string{`foo`, `bar`}},
			}

			for _, cmtTest := range cmtTests {
				Expect(ExtractTokens(cmtTest.str)).To(Equal(cmtTest.tokens))
			}

		})

		It("should extract SGML tokens", func() {

			sgmlTests := []struct {
				str    string
				tokens []string
			}{
				{"<html></html>", []string{"<html>", "</html>"}},
				{"<div id></div>", []string{"<div>", "id", "</div>"}},
				{"<div id=foo></div>", []string{"<div>", "id=", "</div>"}},
				{"<div id class></div>", []string{"<div>", "id", "class", "</div>"}},
				{"<div id=\"foo bar\"></div>", []string{"<div>", "id=", "</div>"}},
				{"<div id='foo bar'></div>", []string{"<div>", "id=", "</div>"}},
				{"<?xml version=\"1.0\"?>", []string{"<?xml>", "version="}},
			}

			for _, sgmlTest := range sgmlTests {
				Expect(ExtractTokens(sgmlTest.str)).To(Equal(sgmlTest.tokens))
			}
		})

	})

	Describe("respect language tokens", func() {

		It("should extract C tokens", func() {

			tokenTests := []struct {
				file   string
				tokens []string
			}{
				{"test_samples/C/hello.h", []string{
					`#ifndef`, `HELLO_H`, `#define`,
					`HELLO_H`, `void`, `hello`,
					`(`, `)`, `;`, `#endif`}},

				{"test_samples/C/hello.c", []string{
					`#include`, `<stdio.h>`, `int`,
					`main`, `(`, `)`,
					`{`, `printf`, `(`,
					`)`, `;`, `return`, `;`, `}`}},
			}

			for _, tokenTest := range tokenTests {
				fileContent, err := ioutil.ReadFile(tokenTest.file)
				if err != nil {
					panic(err)
				}
				Expect(ExtractTokens(string(fileContent))).To(Equal(tokenTest.tokens))
			}
		})

		It("should extract C++ tokens", func() {

			tokenTests := []struct {
				file   string
				tokens []string
			}{
				{"test_samples/C++/bar.h", []string{
					`class`, `Bar`, `{`,
					`protected`, `char`,
					`*name`, `;`, `public`,
					`void`, `hello`, `(`, `)`, `;`, `}`}},

				{"test_samples/C++/hello.cpp", []string{
					`#include`, `<iostream>`, `using`,
					`namespace`, `std`, `;`,
					`int`, `main`, `(`, `)`, `{`,
					`cout`, `<<`, `<<`, `endl`, `;`, `}`}},
			}

			for _, tokenTest := range tokenTests {
				fileContent, err := ioutil.ReadFile(tokenTest.file)
				if err != nil {
					panic(err)
				}
				Expect(ExtractTokens(string(fileContent))).To(Equal(tokenTest.tokens))
			}
		})

		It("should extract Objective-C tokens", func() {

			tokenTests := []struct {
				file   string
				tokens []string
			}{
				{"test_samples/Objective-C/Foo.h", []string{
					`#import`, `<Foundation/Foundation.h>`, `@interface`, `Foo`,
					`NSObject`, `{`, `}`, `@end`}},

				{"test_samples/Objective-C/Foo.m", []string{
					`#import`, `@implementation`, `Foo`, `@end`}},

				{"test_samples/Objective-C/hello.m", []string{
					`#import`, `<Cocoa/Cocoa.h>`, `int`, `main`, `(`, `int`, `argc`,
					`char`, `*argv`, `[`, `]`, `)`, `{`, `NSLog`, `(`, `@`, `)`, `;`,
					`return`, `;`, `}`}},
			}

			for _, tokenTest := range tokenTests {
				fileContent, err := ioutil.ReadFile(tokenTest.file)
				if err != nil {
					panic(err)
				}
				Expect(ExtractTokens(string(fileContent))).To(Equal(tokenTest.tokens))
			}
		})

		It("should extract Objective-C tokens", func() {

			tokenTests := []struct {
				file   string
				tokens []string
			}{
				{"test_samples/Objective-C/Foo.h", []string{
					`#import`, `<Foundation/Foundation.h>`, `@interface`, `Foo`,
					`NSObject`, `{`, `}`, `@end`}},

				{"test_samples/Objective-C/Foo.m", []string{
					`#import`, `@implementation`, `Foo`, `@end`}},

				{"test_samples/Objective-C/hello.m", []string{
					`#import`, `<Cocoa/Cocoa.h>`, `int`, `main`, `(`, `int`, `argc`,
					`char`, `*argv`, `[`, `]`, `)`, `{`, `NSLog`, `(`, `@`, `)`, `;`,
					`return`, `;`, `}`}},
			}

			for _, tokenTest := range tokenTests {
				fileContent, err := ioutil.ReadFile(tokenTest.file)
				if err != nil {
					panic(err)
				}
				Expect(ExtractTokens(string(fileContent))).To(Equal(tokenTest.tokens))
			}
		})

		It("should extract JavaScript tokens", func() {

			tokenTests := []struct {
				file   string
				tokens []string
			}{
				{"test_samples/JavaScript/hello.js", []string{
					`(`, `function`, `(`, `)`, `{`, `console.log`, `(`, `)`, `;`, `}`,
					`)`, `.call`, `(`, `this`, `)`, `;`}},
			}

			for _, tokenTest := range tokenTests {
				fileContent, err := ioutil.ReadFile(tokenTest.file)
				if err != nil {
					panic(err)
				}
				Expect(ExtractTokens(string(fileContent))).To(Equal(tokenTest.tokens))
			}
		})

		It("should extract JSON tokens", func() {

			tokenTests := []struct {
				file   string
				tokens []string
			}{
				{"test_samples/JSON/product.json", []string{
					`{`, `[`, `]`, `{`, `}`, `}`}},
			}

			for _, tokenTest := range tokenTests {
				fileContent, err := ioutil.ReadFile(tokenTest.file)
				if err != nil {
					panic(err)
				}
				Expect(ExtractTokens(string(fileContent))).To(Equal(tokenTest.tokens))
			}
		})
		It("should extract Ruby tokens", func() {

			tokenTests := []struct {
				file   string
				tokens []string
			}{
				{"test_samples/Ruby/foo.rb", []string{
					`module`, `Foo`, `end`}},

				{"test_samples/Ruby/Rakefile", []string{
					`task`, `default`, `do`, `puts`, `end`}},
			}

			for _, tokenTest := range tokenTests {
				fileContent, err := ioutil.ReadFile(tokenTest.file)
				if err != nil {
					panic(err)
				}
				Expect(ExtractTokens(string(fileContent))).To(Equal(tokenTest.tokens))
			}
		})

	})

	Describe("extract shebang", func() {
		It("shoould extract shebang string", func() {

			tokenTests := []struct {
				file    string
				shebang string
			}{
				{"test_samples/Shell/sh", "SHEBANG#!sh"},

				{"test_samples/Shell/bash", "SHEBANG#!bash"},

				{"test_samples/Shell/zsh", "SHEBANG#!zsh"},

				{"test_samples/Shell/invalid-shebang.sh", "echo"},

				{"test_samples/Perl/perl", "SHEBANG#!perl"},

				{"test_samples/Python/python", "SHEBANG#!python"},

				{"test_samples/Ruby/ruby", "SHEBANG#!ruby"},

				{"test_samples/Ruby/ruby2", "SHEBANG#!ruby"},

				{"test_samples/JavaScript/js", "SHEBANG#!node"},

				{"test_samples/PHP/php", "SHEBANG#!php"},

				{"test_samples/Erlang/factorial", "SHEBANG#!escript"},
			}

			for _, tokenTest := range tokenTests {
				fileContent, err := ioutil.ReadFile(tokenTest.file)
				if err != nil {
					panic(err)
				}
				Expect(ExtractTokens(string(fileContent))[0]).To(Equal(tokenTest.shebang))
			}
		})
	})

	Describe("Benchmark", func() {
		fileContent, err := ioutil.ReadFile("./samples/ABAP/24-game-solve.abap")
		if err != nil {
			panic(err)
		}
		Measure("it should tokenizer efficiently", func(b Benchmarker) {
			runtime := b.Time("ExtractTokens runtime", func() {
				ExtractTokens(string(fileContent))
			})

			Expect(runtime.Seconds()).To(BeNumerically("<", 0.05))
		}, 5)
	})
})
