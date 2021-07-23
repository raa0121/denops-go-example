import { Denops, ensureArray, execute, fromFileUrl } from "./deps.ts";

import "./wasm_exec.js";
// deno-lint-ignore no-explicit-any
declare const global: any;

export async function main(denops: Denops) {
  const go = new global.Go();
  const url = fromFileUrl(new URL("main.wasm", import.meta.url));
  const f = await Deno.readFile(url);
  const inst = await WebAssembly.instantiate(f, go.importObject);
  go.run(inst.instance);
  denops.dispatcher = {
    // deno-lint-ignore require-await
    async add(...args: unknown[]) : Promise<void> {
      ensureArray(args);
      let result;
      const nums = (args as string[]).map((i) => parseFloat(i));
      if (nums.some(isNaN)) {
        result = global.GoAdd(args[0], args[1]);
      } else {
        result = global.GoAdd(nums[0], nums[1]);
      }
      console.log(result);
    },
    // deno-lint-ignore require-await
    async addIntNoWrap(...args: unknown[]) : Promise<void> {
      ensureArray(args);
      const nums = (args as string[]).map((i) => parseInt(i));
      const result = global.GoAddIntNoWrap(nums[0], nums[1]);
      console.log(result);
    },
    // deno-lint-ignore require-await
    async addInt(...args: unknown[]) : Promise<void> {
      ensureArray(args);
      const nums = (args as string[]).map((i) => parseInt(i));
      const result = global.GoAddInt(nums[0], nums[1]);
      console.log(result);
    },
    // deno-lint-ignore require-await
    async addFloat(...args: unknown[]) : Promise<void> {
      ensureArray(args);
      const nums = (args as string[]).map((i) => parseFloat(i));
      const result = global.GoAddFloat(nums[0], nums[1]);
      console.log(result);
    },
    // deno-lint-ignore require-await
    async addString(...args: unknown[]) : Promise<void> {
      ensureArray(args);
      const result = global.GoAddString(args[0], args[1]);
      console.log(result);
    },
  };
  await execute(
    denops,
    `command! -nargs=+ GoAdd call denops#request('${denops.name}', 'add', [<f-args>])
     command! -nargs=+ GoAddIntNoWrap call denops#request('${denops.name}', 'addIntNoWrap', [<f-args>])
     command! -nargs=+ GoAddInt call denops#request('${denops.name}', 'addInt', [<f-args>])
     command! -nargs=+ GoAddFloat call denops#request('${denops.name}', 'addFloat', [<f-args>])
     command! -nargs=+ GoAddString call denops#request('${denops.name}', 'addString', [<f-args>])
    `,
  )
}
