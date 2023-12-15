import subprocess
import time
import pandas as pd
import shutil
import os
import matplotlib.pyplot as plt
import numpy as np 
from scipy.interpolate import make_interp_spline
import csv
import openpyxl


# 文件路径
file_path = '/Users/chenchuwen/SOMAS/SOMAS2023Group8/experiment_data/Final_Performance_42.xlsx'

# 使用pandas读取文件
df = pd.read_excel(file_path)
title = "Final Performance"
print(df)
rounds = np.arange(1, 11)  # Rounds 1 to 10
            

# Define colors for the plot lines
grid_color = (0.9, 0.9, 0.9)
lighter_gray_color = (0.8, 0.8, 0.8)
colors = ['salmon','skyblue',lighter_gray_color]

# Lifetime plot
fig, axs = plt.subplots(1, 2, figsize=(18, 5))

for ax in axs:
    ax.grid(True, color=grid_color, linestyle='-', linewidth=0.5)
fig.suptitle(title, fontsize=16)
axs[0].plot(rounds, df["Average_Lifetimes"].tolist(), label='Agent 8', marker='o', linestyle='-', color=colors[0], lw=3,markersize=8)
axs[0].plot(rounds, df["Averages_Lifetimes_Base"].tolist(), label='BaseBiker',marker='o', linestyle='-', color=colors[1], lw=3,markersize=8)
axs[0].plot(rounds, df["Averages_Lifetimes_Other"].tolist(),label='All Other Groups',linestyle='-', color=colors[2], lw=3,markersize=8)
axs[0].set_title('Lifetime Comparison')
axs[0].set_xlabel('Round')
axs[0].set_ylabel('Lifetime')
axs[0].set_ylim(0, 120)
axs[0].set_xlim(1, 10)
axs[0].set_yticks(np.arange(0, 121, 10))
axs[0].set_xticks(np.arange(1, 11))
axs[0].legend(loc='upper right',fontsize=7.5)
axs[0].grid(True)

# Points plot
axs[1].plot(rounds, df["Average_Points"].tolist(), label='Agent 8', marker='o', linestyle='-', color=colors[0], lw=3,markersize=8)
axs[1].plot(rounds, df["Averages_Points_Base"].tolist(), label='BaseBiker',marker='o', linestyle='-', color=colors[1], lw=3,markersize=8)
axs[1].plot(rounds, df["Averages_Points_Other"].tolist(),label='All Other Groups',linestyle='-', color=colors[2], lw=3,markersize=8)
axs[1].set_title('Points Comparison')
axs[1].set_xlabel('Round')
axs[1].set_ylabel('Points')
axs[1].set_xlim(1, 10)
axs[1].set_ylim(0, 12)
axs[1].set_yticks(np.arange(0, 13, 1))
axs[1].set_xticks(np.arange(1, 11))
axs[1].legend(loc='upper right',fontsize=7.5)
axs[1].grid(True)


# Energy plot
# axs[2].plot(rounds, energy_group8, label='Agent 8', marker='o', linestyle='-', color=colors[0], lw=3,markersize=8)
# axs[2].plot(rounds, energy_group0, label='BaseBiker', marker='o',linestyle='-', color=colors[1], lw=3,markersize=8)
# axs[2].plot(rounds, energy_other_groups,label='All Other Groups',linestyle='-', color=colors[2], lw=3,markersize=8)
# axs[2].set_title('Energy Comparison')
# axs[2].set_xlabel('Round')
# axs[2].set_ylabel('Energy')
# axs[2].set_xlim(1, 10)
# axs[2].set_ylim(0, 1.2)
# axs[2].set_yticks(np.arange(0, 1.3, 0.1))
# axs[2].set_xticks(np.arange(1, 11))
# axs[2].legend(loc='upper right',fontsize=7.5)
# axs[2].grid(True)
plt.tight_layout()
fig.savefig(os.path.join('experiment_image', f'{title}.jpg'), dpi=800)