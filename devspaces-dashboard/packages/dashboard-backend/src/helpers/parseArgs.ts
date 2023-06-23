/*
 * Copyright (c) 2018-2023 Red Hat, Inc.
 * This program and the accompanying materials are made
 * available under the terms of the Eclipse Public License 2.0
 * which is available at https://www.eclipse.org/legal/epl-2.0/
 *
 * SPDX-License-Identifier: EPL-2.0
 *
 * Contributors:
 *   Red Hat, Inc. - initial API and implementation
 */

import args from 'args';

export type AppArgs = {
  publicFolder: string;
};

export default function parseArgs(): AppArgs {
  args.option('publicFolder', 'The public folder to serve', './public');

  return args.parse(process.argv) as AppArgs;
}