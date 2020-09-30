#!/usr/bin/env nix-shell
#!nix-shell -i python3 -p python3.pkgs.matplotlib -p python3.pkgs.numpy
import numpy as np
import matplotlib.pyplot as plt
import json
import os
import pprint

pp = pprint.PrettyPrinter(indent=4)

res_f = open("results.json", "r")
res_s = res_f.read()
res_f.close()
s = res_s.split("\n")
res_s_p = ""
for i in range(len(s)):
    res_s_p += s[i] 
    if i < len(s)-2: # FIXME
        res_s_p += ",\n"

res_s_p = "[" + res_s_p + "]"
res_j = json.loads(res_s_p)

idxs = []
for e in ["node-wot", "wot-py"]:
    with open("tests/" + e + "/index.json", "r") as f:
        idxs += json.load(f)

try:
    os.mkdir("plots")
except FileExistsError:
    pass

for t in res_j:
    if t["type"] != 1:
        continue
    name = t["name"]

    # FIXME
    spec = list(filter(lambda x: x["path"] == t["path"] \
            # and set(x["args"]) == set(t["args"]) \
            and x["type"] == t["type"], idxs))
    pp.pprint(spec)
    idxs.remove(spec[0])

    x = []
    xlabels = []
    i = 0
    for m in t["measurements"]:
        x.append(np.array([]))
        not_so_raw = filter(lambda x: int(x["duration"]) >= 0, t["measurements"][m]["raw"])
        for r in not_so_raw:
            x[i] = np.append(x[i], np.float64(r["duration"])/1000)
        x[i] = x[i].astype(np.float64)
        rs = spec[0]["measureTestProperties"]["requestSets"][i]
        xlabels.append(str(rs["num"])+ "\n@" + str(rs["parallel"]))
        i+=1

    fig1, ax1 = plt.subplots(figsize=(7, 5))
    fig1.set_dpi(200)
    ax1.set_title(name)
    ax1.set_ylabel('Anfragebearbeitungszeit in ms')
    ax1.set_xlabel('Anfragen gesamt @ Anzahl paralleler Verbindungen')

    ax1.boxplot(x, labels=xlabels)
    fig1.tight_layout()
    fig1.savefig("plots/"+name+".pdf")
