# Use a specific, stable Python version as the base image
FROM python:3.12-slim

# Set environment variables to prevent generating .pyc files and to run python in unbuffered mode
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

# Set the working directory inside the container
WORKDIR /app

# The application code is in the 'server' subdirectory.
# Copy the requirements file first to leverage Docker's layer caching.
COPY server/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of the application code from the 'server' directory into the container
COPY server/ .

# Expose the port the application will run on
EXPOSE 8000

# Command to run the application using uvicorn when the container starts
CMD ["uvicorn", "app_distribution_server.app:app", "--host", "0.0.0.0", "--port", "8000"]
