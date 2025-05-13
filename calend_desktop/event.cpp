#include "event.h"
#include "ui_event.h"

Event::Event(QString uid, QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::Event), cli(settings.value("host").toString(), settings.value("port").toString()), uid(uid)
{
    ui->setupUi(this);
}

Event::~Event()
{
    delete ui;

}


void Event::on_create_event_clicked()
{
    ui->create_event->setDisabled(true);
    EventData event;
    event.master = uid;
    event.name = ui->nameEdit->text();
    event.desc = ui->descEdit->text();
    event.type = ui->typeCombo->currentText();
    event.prior = static_cast<EventData::Priority>(ui->priorityCombo->currentIndex() + 1);

    QDate startDate = ui->startDateEdit->date();
    QTime startTime = ui->startTimeEdit->time();
    event.start = QDateTime(startDate, startTime);


    QDate endDate = ui->endDateEdit->date();
    QTime endTime = ui->endTimeEdit->time();
    event.end = QDateTime(endDate, endTime);


    QStringList parts = ui->parts->toPlainText().split(" ");
    event.parts = parts;


    bool ok = cli.addEvent(event, uid);

    if (!ok) {
        ui->result->setText("Ошибка создания события");
        ui->create_event->setEnabled(true);
    } else {
        this->close();
    }
}

