import pandas as pd
import matplotlib.pyplot as plt
from matplotlib.animation import FuncAnimation


def load_data(file_path):
    """ Load simulation data from a CSV file. """
    try:
        return pd.read_csv(file_path)
    except Exception as e:
        print(f"Failed to read the file: {e}")
        return None


def update(frame_number, data, scatters):
    """ Update function for animation, sets the data for each scatter plot. """
    for scatter, body in zip(scatters, data['Body Name'].unique()):
        body_data = data[(data['Frame'] == frame_number) & (data['Body Name'] == body)]
        if not body_data.empty:
            pos_data = body_data[['PosX', 'PosY']].values
            if pos_data.size > 0:
                scatter.set_offsets(pos_data)
            else:
                print(f"No position data available for {body} at frame {frame_number}")
        else:
            print(f"No data for {body} at frame {frame_number}")
    return scatters


def animate_trajectories(df):
    """ Create an animation of the trajectories of each body. """
    if df is None or df.empty:
        print("No data available to plot.")
        return

    # Setup the plot
    fig, ax = plt.subplots()
    ax.set_xlim(df['PosX'].min() - 100, df['PosX'].max() + 100)
    ax.set_ylim(df['PosY'].min() - 100, df['PosY'].max() + 100)

    # Create a scatter plot for each body
    scatters = [ax.scatter(df['PosX'], df['PosY'], label=body) for body in df['Body Name'].unique()]

    plt.title('Simulation of Body Movements')
    plt.xlabel('Position X')
    plt.ylabel('Position Y')
    plt.legend()
    plt.grid(True)

    # Create the animation
    frames = df['Frame'].max() + 1
    ani = FuncAnimation(fig, update, frames=frames, fargs=(df, scatters), interval=50, blit=True)

    # Save the animation
    ani.save('body_movements.mp4', writer='ffmpeg', fps=60)


def main():
    # Path to the CSV file
    file_path = 'simulation_results.csv'

    # Load the data
    df = load_data(file_path)

    # Animate trajectories
    animate_trajectories(df)


if __name__ == "__main__":
    main()
