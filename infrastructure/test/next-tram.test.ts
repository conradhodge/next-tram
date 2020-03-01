import { expect as expectCDK, matchTemplate, MatchStyle } from '@aws-cdk/assert';
import * as cdk from '@aws-cdk/core';
import NextTram = require('../lib/next-tram-stack');

test('Empty Stack', () => {
    const app = new cdk.App();
    // WHEN
    const stack = new NextTram.NextTramStack(app, 'MyTestStack');
    // THEN
    expectCDK(stack).to(matchTemplate({
      "Resources": {}
    }, MatchStyle.EXACT))
});
