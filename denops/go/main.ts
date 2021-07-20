import { Denops, ensureArray, execute, fromFileUrl } from "./deps.ts";

import "./wasm_exec.js";
declare const global: any;

export async function main(denops: Denops) {
  const go = new global.Go();
  const url = fromFileUrl(new URL("main.wasm", import.meta.url));
  const f = await Deno.readFile(url);
  const inst = await WebAssembly.instantiate(f, go.importObject);
  go.run(inst.instance);
  denops.dispatcher = {
    async add(...args: unknown[]) : Promise<void> {
      ensureArray(args);
      const nums = (args as string[]).map((i) => parseInt(i));
      const result = global.GoAdd(nums[0], nums[1]);
      console.log(result);
    },
  };
  await execute(
    denops,
    `command! -nargs=+ GoAdd call denops#request('${denops.name}', 'add', [<f-args>])`,
  )
}
