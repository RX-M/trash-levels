namespace * Waste

/**
 * Service for collecting real time IoT trash level information.
 */
service TrashLevels {
    /**
     * Return the percentage of trash can capacity consumed for a specified trash can
     *
     * @param the trash can ID
     * @return the percentage full
     */
    i16 GetLevel( 1: i32 canid )
}
