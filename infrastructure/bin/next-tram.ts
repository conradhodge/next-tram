#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { NextTramStack } from '../lib/next-tram-stack';

const app = new cdk.App();
new NextTramStack(app, 'next-tram-stack');
