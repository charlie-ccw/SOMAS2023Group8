import subprocess
import time
import pandas as pd
import shutil
import os
import matplotlib.pyplot as plt
import numpy as np 
from scipy.interpolate import make_interp_spline

# data = {
#     'Average_Lifetimes': [],
#     'Average_Points': [],
#     'Averages_Lifetimes_Base': [],
#     'Averages_Points_Base': []
# }

# df = pd.DataFrame(data)

# # Save the empty DataFrame to a CSV file
# csv_file = 'kickout.csv'  # Replace with the desired file name
# df.to_csv(csv_file, index=False)

# Rounds 1 to 10
# Specify the path to your CSV file
csv_file = 'decidejoin_thres.csv'  # Replace with the path to your CSV file

# Read the CSV file into a DataFrame
df = pd.read_csv(csv_file)

# Extract the data for Lifetime and Points into separate sets
lifetime_set1 = df['Average_Lifetimes'].tolist()
lifetime_set2 = df['Averages_Lifetimes_Base'].tolist()
points_set1 = df['Average_Points'].tolist()
points_set2 = df['Averages_Points_Base'].tolist()
highlight_x_value = 0.2
rounds = [0.5,0.4,0.3,0.2,0.1,0.0,-0.1,-0.2]
# Average_lifetime = [52.85,57.375,64.1,50.025,56.85,39.675,53.375,48.775]
# Average_points = [3.072196073237026,3.090323781646945,3.6158923076081293,1.1900832497151286,2.2899678313857494,1.1904207018989834,1.4441601731235632,2.557598534762108]
# Average_energy = [0.6479229460437743,0.6828458913961912,0.6939616409729827,0.6482222304828226,0.6670848178502089,0.6511782315631169,0.6533695843814554,0.628460462325583]
# Define colors for the plot lines
grid_color = (0.9, 0.9, 0.9)
lighter_gray_color = (0.8, 0.8, 0.8)
colors = ['red','skyblue','salmon','#1f77b4']

# Lifetime plot
# fig, axs = plt.subplots(1, 3, figsize=(18, 5))

# for ax in axs:
#     ax.grid(True, color=grid_color, linestyle='-', linewidth=0.5)
    
title = "Thershold of socre about accpeting agent for Agent Recruitment"
# fig.suptitle(title, fontsize=16)
# axs[0].plot(rounds, Average_lifetime, label='Agent 8', marker='o', linestyle='-', color=colors[0], lw=3,markersize=8)
# axs[0].set_title('Lifetime Comparison')
# axs[0].set_xlabel('Kickout Threshold')
# axs[0].set_ylabel('Lifetime')
# axs[0].set_ylim(0, 120)
# axs[0].set_xlim(0.0, -0.8,-0.1)
# axs[0].set_yticks(np.arange(0, 121, 10))
# axs[0].set_xticks(np.arange(0.0, -0.8,-0.1))
# axs[0].legend(loc='upper right',fontsize=7.5)
# axs[0].grid(True)

# # Points plot
# axs[1].plot(rounds, Average_points, label='Agent 8', marker='o', linestyle='-', color=colors[0], lw=3,markersize=8)
# # axs[1].plot(rounds, points_group0, label='BaseBiker',marker='o', linestyle='-', color=colors[1], lw=3,markersize=8)
# # axs[1].plot(rounds, points_other_groups,label='All Other Groups',linestyle='-', color=colors[2], lw=3,markersize=8)
# axs[1].set_title('Points Comparison')
# axs[1].set_xlabel('Kickout Threshold')
# axs[1].set_ylabel('Points')
# axs[1].set_xlim(0.0, -0.8,-0.1)
# axs[1].set_ylim(0, 12)
# axs[1].set_yticks(np.arange(0, 13, 1))
# axs[1].set_xticks(np.arange(0.0, -0.8,-0.1))
# axs[1].legend(loc='upper right',fontsize=7.5)
# axs[1].grid(True)


# # Energy plot
# axs[2].plot(rounds, Average_energy, label='Agent 8', marker='o', linestyle='-', color=colors[0], lw=3,markersize=8)
# # axs[2].plot(rounds, energy_group0, label='BaseBiker', marker='o',linestyle='-', color=colors[1], lw=3,markersize=8)
# # axs[2].plot(rounds, energy_other_groups,label='All Other Groups',linestyle='-', color=colors[2], lw=3,markersize=8)
# axs[2].set_title('Energy Comparison')
# axs[2].set_xlabel('Kickout Threshold')
# axs[2].set_ylabel('Energy')
# axs[2].set_xlim(0.0, -0.8,-0.1)
# axs[2].set_ylim(0, 1.2)
# axs[2].set_yticks(np.arange(0, 1.3, 0.1))
# axs[2].set_xticks(np.arange(0.0, -0.8,-0.1))
# axs[2].legend(loc='upper right',fontsize=7.5)
# axs[2].grid(True)
# plt.tight_layout()


# fig, ax = plt.subplots(figsize=(8, 5))  # Adjust the figsize as needed

# # Define colors for the plot lines and data points
# colors = ['salmon', 'skyblue', 'lightgray']  # You can add more colors if needed

# # Plot the lines
# ax.plot(rounds, Average_lifetime, label='Lifetime', marker='o', linestyle='-', color=colors[0], lw=3, markersize=8)
# ax.plot(rounds, Average_points, label='Points', marker='o', linestyle='-', color=colors[1], lw=3, markersize=8)
# ax.plot(rounds, Average_energy, label='Energy', marker='o',linestyle='-', color=colors[2], lw=3, markersize=8)

# # Plot the data points as scatter plots
# ax.annotate(f'{Average_lifetime[2]:.2f}', (rounds[2], Average_lifetime[3]), textcoords="offset points", xytext=(0, 10), ha='center')
# ax.annotate(f'{Average_points[2]:.2f}', (rounds[2], Average_points[3]), textcoords="offset points", xytext=(0, 10), ha='center')
# ax.annotate(f'{Average_energy[2]:.2f}', (rounds[2], Average_energy[3]), textcoords="offset points", xytext=(0, 10), ha='center')

# # Set the title and labels
# ax.set_title(title)
# ax.set_xlabel('Kickout Threshold')

# # Set the x and y limits and ticks
# ax.set_xlim(0.0, -0.8)
# ax.set_xticks(np.arange(0.0, -0.8, -0.1))

# # Add a legend
# ax.legend(loc='upper right', fontsize=7.5)

# # Add a grid
# ax.grid(True, color=(0.8, 0.8, 0.8), linestyle='--', linewidth=0.5)  # Use a lighter gray for the grid



# fig.savefig(os.path.join('experiment_image', f'{title}.jpg'), dpi=300)  # Specify the folder path and desired file name
# plt.show()

fig, ax1 = plt.subplots(figsize=(8, 5))
fig.suptitle(title, fontsize=16)
# Plot the lines for the first dataset on the primary y-axis
ax1.plot(rounds, points_set1, label='Points_Agent 8 ', marker='o', linestyle='-', color=colors[0], lw=3, markersize=8)
ax1.plot(rounds, points_set2, label='Points_BaseBiker', marker='o', linestyle='-', color=colors[1], lw=3, markersize=8)

# Set the title and labels for the primary y-axis
ax1.set_title('Points and life time comparison with BaseAgent',fontsize = 10)
ax1.set_xlabel('Threshold')
ax1.set_ylabel('Points', color='black')  # Label for the primary y-axis

# Set the limits and ticks for the primary y-axis
ax1.set_ylim(0, 12)  # Adjust the limits as needed
ax1.set_yticks(np.arange(0, 10.0, 1.0))

# Create a secondary y-axis for the third dataset
ax2 = ax1.twinx()

# Plot the lines for the third dataset on the secondary y-axis
ax2.plot(rounds, lifetime_set1, label='Lifetime_Agent 8', marker ='o',linestyle='-', color=colors[2], lw=3, markersize=8)
ax2.plot(rounds, lifetime_set2, label='Lifetime_Basebiker',marker ='o', linestyle='-', color=colors[3], lw=3, markersize=8)
ax2.axvline(x=highlight_x_value, color='#dc3912', linestyle='--', lw=2,label = 'Agent 8 Threshold')
# Label the data points for the third dataset with numerical values

# Set the labels for the secondary y-axis
ax2.set_ylabel('Lifetime', color='black')  # Label for the secondary y-axis

# Set the limits and ticks for the secondary y-axis
ax2.set_ylim(0, 130)  # Adjust the limits as needed
ax2.set_yticks(np.arange(0, 101, 10))

# Add a legend
lines, labels = ax1.get_legend_handles_labels()
lines2, labels2 = ax2.get_legend_handles_labels()
ax1.legend(lines + lines2, labels + labels2, loc='upper right', fontsize=7.5)

# Add a grid
ax1.grid(True, color=(0.8, 0.8, 0.8), linestyle='--', linewidth=0.5)  # Use a lighter gray for the grid

# Save the figure
fig.savefig(os.path.join('experiment_image', f'{title}.jpg'), dpi=300)

# Display the plot
plt.show()